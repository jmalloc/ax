package saga

import (
	"context"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/internal/typeswitch"
)

// Workflow is a Saga for implementing application-defined workflows.
type Workflow struct {
	ErrorIfNotFound
	CompletableByData

	Prototype   Data
	Triggers    ax.MessageTypeSet
	NonTriggers ax.MessageTypeSet
	Handle      typeswitch.Switch
}

// NewWorkflow returns a saga that forwards to the given aggregate.
//
// Workflows are a specialization of sagas that handle events and produce
// commands.
//
// It accepts a prototype data instance which is cloned for new instances.
//
// For each event type to be handled, the aggregate must implement a "handler"
// method that adheres to one of the following signatures:
//
//     func (ev *<T>) []ax.Command
//     func (ev *<T>, env ax.Envelope) []ax.Command
//
// Where T is a struct type that implements ax.Event.
//
// Handler methods are responsible for mutating the state of the workflow and
// producing new commands, based on the event being handled. They may inspect
// the current state of the workflow, and then return zero or more commands
// to be executed.
//
// The names of handler methods are meaningful to the workflow system. If an
// event is meant to trigger a new workflow instance, its handler method's name
// must begin with "StartWhen". Other handler methods must begin with "When". By
// convention these prefixes are followed by the message name, such as:
//
//     func (*BankTransferWorkflow) StartWhenTransferStarted(*messages.TransferStarted)
//     func (*BankTransferWorkflow) WhenAccountDebited(*messages.AccountDebited)
func NewWorkflow(p Data) *Workflow {
	w := &Workflow{
		Prototype: p,
	}

	// setup type-switch for event handlers.
	sw, types, err := typeswitch.New(
		[]reflect.Type{
			reflect.TypeOf(p),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
		},
		[]reflect.Type{
			reflect.TypeOf(([]ax.Command)(nil)),
		},
		workflowStartWhenSignature,
		workflowStartWhenWithEnvelopeSignature,
		workflowWhenSignature,
		workflowWhenWithEnvelopeSignature,
	)
	if err != nil {
		panic(err)
	}

	w.Handle = sw

	w.Triggers = ax.TypesByGoType(
		append(
			types[workflowStartWhenSignature],
			types[workflowStartWhenWithEnvelopeSignature]...,
		)...,
	)

	w.NonTriggers = ax.TypesByGoType(
		append(
			types[workflowWhenSignature],
			types[workflowWhenWithEnvelopeSignature]...,
		)...,
	)

	return w
}

// PersistenceKey returns a unique identifier for the saga.
//
// The persistence key is used to relate persisted data with the saga
// implementation that owns it. Persistence keys should not be changed once
// a saga has active instances.
func (w *Workflow) PersistenceKey() string {
	return proto.MessageName(w.Prototype)
}

// MessageTypes returns the set of messages that are routed to this saga.
//
// tr is the set of "trigger" messages. If they can not be routed to an
// existing saga instance a new instance is created.
//
// mt is the set of messages that are only routed to existing instances. If
// they can not be routed to an existing instance, the HandleNotFound()
// method is called instead.
func (w *Workflow) MessageTypes() (tr ax.MessageTypeSet, mt ax.MessageTypeSet) {
	return w.Triggers, w.NonTriggers
}

// NewData returns a pointer to a new zero-value instance of the
// saga's data type.
func (w *Workflow) NewData() Data {
	return proto.Clone(w.Prototype).(Data)
}

// HandleMessage handles a message for a particular saga instance.
func (w *Workflow) HandleMessage(ctx context.Context, s ax.Sender, env ax.Envelope, i Instance) error {
	out := w.Handle.Dispatch(
		i.Data,
		env.Message.(ax.Event),
		env,
	)

	if cmds := out[0]; cmds != nil {
		for _, m := range cmds.([]ax.Command) {
			if _, err := s.ExecuteCommand(ctx, m); err != nil {
				return err
			}
		}
	}

	return nil
}

var (
	workflowStartWhenSignature = &typeswitch.Signature{
		Prefix: "StartWhen",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
		},
		Out: []reflect.Type{
			reflect.TypeOf(([]ax.Command)(nil)),
		},
	}

	workflowStartWhenWithEnvelopeSignature = &typeswitch.Signature{
		Prefix: "StartWhen",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
			reflect.TypeOf((*ax.Envelope)(nil)).Elem(),
		},
		Out: []reflect.Type{
			reflect.TypeOf(([]ax.Command)(nil)),
		},
	}

	workflowWhenSignature = &typeswitch.Signature{
		Prefix: "When",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
		},
		Out: []reflect.Type{
			reflect.TypeOf(([]ax.Command)(nil)),
		},
	}

	workflowWhenWithEnvelopeSignature = &typeswitch.Signature{
		Prefix: "When",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
			reflect.TypeOf((*ax.Envelope)(nil)).Elem(),
		},
		Out: []reflect.Type{
			reflect.TypeOf(([]ax.Command)(nil)),
		},
	}
)
