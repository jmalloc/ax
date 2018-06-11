package aggregate

import (
	"context"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/internal/visitor"
)

// Saga is an implementation of saga.Saga that wraps an Aggregate.
type Saga struct {
	saga.ErrorIfNotFound
	saga.CompletableByData

	Prototype  Aggregate
	Identifier Identifier

	Triggers ax.MessageTypeSet
	Handle   func(Aggregate, ax.Command, Recorder)
	Apply    func(Aggregate, ax.Event)
}

// New returns a saga that forwards to the given aggregate.
func New(agg Aggregate, opts ...Option) *Saga {
	sg := &Saga{
		Prototype: agg,
	}

	for _, opt := range opts {
		opt(sg)
	}

	commandTypes := visitor.MakeAcceptor(
		&sg.Handle,
		reflect.TypeOf((*ax.Command)(nil)).Elem(),
		reflect.TypeOf(agg),
	)

	visitor.MakeAcceptor(
		&sg.Apply,
		reflect.TypeOf((*ax.Event)(nil)).Elem(),
		reflect.TypeOf(agg),
	)

	// TODO: make use of https://github.com/jmalloc/ax/issues/74
	for _, t := range commandTypes {
		sg.Triggers = sg.Triggers.Add(
			ax.TypeOf(
				reflect.Zero(t).Interface().(ax.Message),
			),
		)
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
	return sg.Triggers, ax.MessageTypeSet{}
}

// GenerateInstanceID returns the saga ID to use for a new instance.
//
// It is called when a "trigger" message is received and there is no
// existing saga instance. env contains the "trigger" message.
func (sg *Saga) GenerateInstanceID(ctx context.Context, env ax.Envelope) (saga.InstanceID, error) {
	id, err := sg.Identifier.AggregateID(
		env.Message.(ax.Command),
	)

	return saga.InstanceID{ID: id.ID}, err
}

// NewData returns a pointer to a new zero-value instance of the
// saga's data type.
func (sg *Saga) NewData() saga.Data {
	return proto.Clone(sg.Prototype).(saga.Data)
}

// InstanceIDForMessage returns the ID of the saga instance to which the
// given message is routed, if any.
//
// If ok is false the message is ignored; otherwise, the message is routed
// to the saga instance with the returned ID.
func (sg *Saga) InstanceIDForMessage(ctx context.Context, env ax.Envelope) (saga.InstanceID, bool, error) {
	id, err := sg.Identifier.AggregateID(
		env.Message.(ax.Command),
	)

	return saga.InstanceID{ID: id.ID}, true, err
}

// HandleMessage handles a message for a particular saga instance.
func (sg *Saga) HandleMessage(ctx context.Context, s ax.Sender, env ax.Envelope, i saga.Instance) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	rec := func(m ax.Event) {
		if _, err := s.PublishEvent(ctx, m); err != nil {
			panic(err)
		}
	}

	sg.Handle(
		i.Data,
		env.Message.(ax.Command),
		rec,
	)

	return
}

// ApplyEvent updates d to reflect the fact that an event has occurred.
//
// It may panic if env.Message does not implement ax.Event.
func (sg *Saga) ApplyEvent(d saga.Data, env ax.Envelope) {
	m := env.Message.(ax.Event)
	sg.Apply(d, m)
}
