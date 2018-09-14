package saga

import (
	"context"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/internal/typeswitch"
)

// Workflow is a Saga for implementing application-defined workflows.
type Workflow struct {
	ErrorIfNotFound
	CompletableByData

	Prototype     Data
	Triggers      ax.MessageTypeSet
	NonTriggers   ax.MessageTypeSet
	CommandHandle typeswitch.Switch
	EventHandle   typeswitch.Switch
}

// NewWorkflow returns a saga that forwards to the given aggregate.
//
// Workflows are a specialization of sagas that handle commands and/or events
// and produce commands.
//
// It accepts a prototype data instance which is cloned for new instances.
//
// For each message type to be handled, the aggregate must implement a "handler"
// method that adheres to one of the following signatures:
//
//     func (m *<T>, ax.CommandExecutor)
//     func (m *<T>, env ax.Envelope, ax.CommandExecutor)
//
// Where T is a struct type that implements ax.Message.
//
// Handler methods are responsible for mutating the state of the workflow and
// producing new commands, based on the message being handled.
//
// The names of handler methods are meaningful to the workflow system. If a
// command or event is meant to trigger a new workflow instance, its handler
// method's name must have a prefix "Begin". Non-triggering command handler
// methods must have prefix "Do". Both triggering and non-triggering event
// handler methods must have prefix "When". For triggering event handlers the
// combined prefix results in "BeginWhen". By convention these prefixes are
// followed by the message name, such as:
//
//      // workflow-triggering command handler
//      func (*BankTransferWorkflow) BeginDebitAccount(
//              *messages.DebitAccount,
//              ax.CommandExecutor,
//          )
//
//      // non-triggering command handler
//      func (*BankTransferWorkflow) DoDebitAccount(
//              *messages.DebitAccount,
//              ax.CommandExecutor,
//      )
//
//      // workflow-triggering event handler
//      func (*BankTransferWorkflow) BeginWhenAccountDebited(
//              *messages.AccountDebited,
//              ax.CommandExecutor,
//      )
//
//      // non-triggering event handler
//      func (*BankTransferWorkflow) WhenAccountDebited(
//              *messages.AccountDebited,
//              ax.CommandExecutor,
//      )
func NewWorkflow(p Data) *Workflow {
	w := &Workflow{
		Prototype: p,
	}

	// setup type-switch for command handlers.
	csw, ctypes, err := typeswitch.New(
		[]reflect.Type{
			reflect.TypeOf(p),
			reflect.TypeOf((*ax.Command)(nil)).Elem(),
			reflect.TypeOf((*ax.Envelope)(nil)).Elem(),
			reflect.TypeOf((*ax.CommandExecutor)(nil)).Elem(),
		},
		nil,
		workflowBeginSignature,
		workflowBeginWithEnvelopeSignature,
		workflowDoSignature,
		workflowDoWithEnvelopeSignature,
	)
	if err != nil {
		panic(err)
	}

	w.CommandHandle = csw

	// setup type-switch for event handlers.
	esw, etypes, err := typeswitch.New(
		[]reflect.Type{
			reflect.TypeOf(p),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
			reflect.TypeOf((*ax.Envelope)(nil)).Elem(),
			reflect.TypeOf((*ax.CommandExecutor)(nil)).Elem(),
		},
		nil,
		workflowBeginWhenSignature,
		workflowBeginWhenWithEnvelopeSignature,
		workflowWhenSignature,
		workflowWhenWithEnvelopeSignature,
	)
	if err != nil {
		panic(err)
	}

	w.EventHandle = esw

	w.Triggers = ax.TypesByGoType(
		mergeTypeSlices(
			ctypes[workflowBeginSignature],
			ctypes[workflowBeginWithEnvelopeSignature],
			etypes[workflowBeginWhenSignature],
			etypes[workflowBeginWhenWithEnvelopeSignature],
		)...,
	)

	w.NonTriggers = ax.TypesByGoType(
		mergeTypeSlices(
			ctypes[workflowDoSignature],
			ctypes[workflowDoWithEnvelopeSignature],
			etypes[workflowWhenSignature],
			etypes[workflowWhenWithEnvelopeSignature],
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
func (w *Workflow) HandleMessage(
	ctx context.Context,
	s ax.Sender,
	env ax.Envelope,
	i Instance,
) error {
	type command struct {
		Command ax.Command
		Options []ax.ExecuteOption
	}

	var cmds []command

	switch t := env.Message.(type) {
	case ax.Command:
		w.CommandHandle.Dispatch(
			i.Data,
			env.Message.(ax.Command),
			env,
			func(m ax.Command, opts ...ax.ExecuteOption) {
				cmds = append(cmds, command{m, opts})
			},
		)
	case ax.Event:
		w.EventHandle.Dispatch(
			i.Data,
			env.Message.(ax.Event),
			env,
			func(m ax.Command, opts ...ax.ExecuteOption) {
				cmds = append(cmds, command{m, opts})
			},
		)
	default:
		return fmt.Errorf(
			"unknown message type %T",
			t,
		)
	}

	for _, c := range cmds {
		if _, err := s.ExecuteCommand(ctx, c.Command, c.Options...); err != nil {
			return err
		}
	}

	return nil
}

var (
	workflowBeginSignature = &typeswitch.Signature{
		Prefix: "Begin",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Command)(nil)).Elem(),
			reflect.TypeOf((*ax.CommandExecutor)(nil)).Elem(),
		},
	}

	workflowBeginWithEnvelopeSignature = &typeswitch.Signature{
		Prefix: "Begin",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Command)(nil)).Elem(),
			reflect.TypeOf((*ax.Envelope)(nil)).Elem(),
			reflect.TypeOf((*ax.CommandExecutor)(nil)).Elem(),
		},
	}

	workflowDoSignature = &typeswitch.Signature{
		Prefix: "Do",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Command)(nil)).Elem(),
			reflect.TypeOf((*ax.CommandExecutor)(nil)).Elem(),
		},
	}

	workflowDoWithEnvelopeSignature = &typeswitch.Signature{
		Prefix: "Do",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Command)(nil)).Elem(),
			reflect.TypeOf((*ax.Envelope)(nil)).Elem(),
			reflect.TypeOf((*ax.CommandExecutor)(nil)).Elem(),
		},
	}

	workflowBeginWhenSignature = &typeswitch.Signature{
		Prefix: "BeginWhen",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
			reflect.TypeOf((*ax.CommandExecutor)(nil)).Elem(),
		},
	}

	workflowBeginWhenWithEnvelopeSignature = &typeswitch.Signature{
		Prefix: "BeginWhen",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
			reflect.TypeOf((*ax.Envelope)(nil)).Elem(),
			reflect.TypeOf((*ax.CommandExecutor)(nil)).Elem(),
		},
	}

	workflowWhenSignature = &typeswitch.Signature{
		Prefix: "When",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
			reflect.TypeOf((*ax.CommandExecutor)(nil)).Elem(),
		},
	}

	workflowWhenWithEnvelopeSignature = &typeswitch.Signature{
		Prefix: "When",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
			reflect.TypeOf((*ax.Envelope)(nil)).Elem(),
			reflect.TypeOf((*ax.CommandExecutor)(nil)).Elem(),
		},
	}
)

// mergeTypeSlices appends all slices of reflect.Type to a single slice.
func mergeTypeSlices(slices ...[]reflect.Type) []reflect.Type {
	var r []reflect.Type
	for _, s := range slices {
		r = append(r, s...)
	}
	return r
}
