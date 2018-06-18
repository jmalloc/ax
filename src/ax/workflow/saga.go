package workflow

import (
	"context"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/internal/visitor"
)

// Saga is an implementation of saga.Saga that wraps a Workflow.
type Saga struct {
	saga.ErrorIfNotFound
	saga.CompletableByData

	Prototype Workflow

	Triggers      ax.MessageTypeSet
	NonTriggers   ax.MessageTypeSet
	HandleTrigger func(Workflow, ax.Event) []ax.Command
	Handle        func(Workflow, ax.Event) []ax.Command
}

// New returns a saga that forwards to the given aggregate.
func New(w Workflow) *Saga {
	sg := &Saga{
		Prototype: w,
	}

	triggerTypes := visitor.MakeAcceptor(
		&sg.HandleTrigger,
		reflect.TypeOf((*ax.Event)(nil)).Elem(),
		reflect.TypeOf(w),
		"StartWhen",
	)

	nonTriggerTypes := visitor.MakeAcceptor(
		&sg.Handle,
		reflect.TypeOf((*ax.Event)(nil)).Elem(),
		reflect.TypeOf(w),
		"When",
	)

	sg.Triggers = ax.TypesByGoType(triggerTypes...)
	sg.NonTriggers = ax.TypesByGoType(nonTriggerTypes...)

	for _, mt := range sg.Triggers.Intersection(sg.NonTriggers).Members() {
		panic("workflow defines multiple handlers for " + mt.Name)
	}

	return sg
}

// PersistenceKey returns a unique identifier for the saga.
//
// The persistence key is used to relate persisted data with the saga
// implementation that owns it. Persistence keys should not be changed once
// a saga has active instances.
func (sg *Saga) PersistenceKey() string {
	return proto.MessageName(sg.Prototype)
}

// MessageTypes returns the set of messages that are routed to this saga.
//
// tr is the set of "trigger" messages. If they can not be routed to an
// existing saga instance a new instance is created.
//
// mt is the set of messages that are only routed to existing instances. If
// they can not be routed to an existing instance, the HandleNotFound()
// method is called instead.
func (sg *Saga) MessageTypes() (tr ax.MessageTypeSet, mt ax.MessageTypeSet) {
	return sg.Triggers, sg.NonTriggers
}

// NewData returns a pointer to a new zero-value instance of the
// saga's data type.
func (sg *Saga) NewData() saga.Data {
	return proto.Clone(sg.Prototype).(saga.Data)
}

// HandleMessage handles a message for a particular saga instance.
func (sg *Saga) HandleMessage(ctx context.Context, s ax.Sender, env ax.Envelope, i saga.Instance) error {
	h := sg.Handle
	if sg.Triggers.Has(env.Type()) {
		h = sg.HandleTrigger
	}

	cmds := h(
		i.Data,
		env.Message.(ax.Event),
	)

	for _, m := range cmds {
		if _, err := s.ExecuteCommand(ctx, m); err != nil {
			return err
		}
	}

	return nil
}
