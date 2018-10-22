package saga

import (
	"context"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/internal/typeswitch"
)

// Aggregate is a Saga for implementing application-defined domain aggregates.
type Aggregate struct {
	ErrorIfNotFound
	CompletableByData

	Prototype Data
	Triggers  ax.MessageTypeSet
	Handle    typeswitch.Switch
	Apply     typeswitch.Switch
}

// NewAggregate returns a new aggregate saga.
//
// Aggregates are a specialization of sagas that handle commands and produce
// events.
//
// It accepts a prototype data instance which is cloned for new instances.
//
// For each command type to be handled, the aggregate's data struct must
// implement a "handler" method that adheres to one of the following signatures:
//
//     func (cmd *<T>, rec ax.EventRecorder)
//     func (cmd *<T>, mctx ax.MessageContext, rec ax.EventRecorder)
//
// Where T is a struct type that implements ax.Command.
//
// Handler methods are responsible for producing new events based on the command
// being handled. They may inspect the current state of the aggregate, and then
// record zero or more events using rec. Handlers should never mutate the
// aggregate state.
//
// The names of handler methods are meaningful. Each handler method's name must
// begin with "Do". By convention these prefixes are followed by the message
// name, such as:
//
//     func (*BankAccount) DoCreditAccount(*messages.CreditAccount, ax.EventRecorder)
//
// For each of the event types passed to rec, the aggregate must implement an
// "applier" method that adheres to one of the following signatures:
//
//     func (ev *T)
//     func (ev *T, mctx ax.MessageContext)
//
// Where T is a struct type that implements ax.Event.
//
// Applier methods are responsible for mutating the aggregate state. The applier
// is called every time an event is recorded, *and* when loading an
// event-sourced aggregate from the message store.
//
// The names of handler methods are meaningful. Each handler method's name must
// begin with "When". By convention these prefixes are followed by the message
// name, such as:
//
//     func (*BankAccount) WhenAccountCredited(*messages.AccountCredited)
func NewAggregate(p Data) *Aggregate {
	a := &Aggregate{
		Prototype: p,
	}

	// setup type-switch for command handlers.
	sw, _, err := typeswitch.New(
		[]reflect.Type{
			reflect.TypeOf(p),
			reflect.TypeOf((*ax.Command)(nil)).Elem(),
			reflect.TypeOf((*ax.MessageContext)(nil)).Elem(),
			reflect.TypeOf((*ax.EventRecorder)(nil)).Elem(),
		},
		nil, // no outputs
		aggregateHandleSignature,
		aggregateHandleSignatureWithMessageContext,
	)
	if err != nil {
		panic(err)
	}

	a.Handle = sw
	a.Triggers = ax.TypesByGoType(sw.Types()...)

	// setup type-switch for event appliers.
	a.Apply, _, err = typeswitch.New(
		[]reflect.Type{
			reflect.TypeOf(p),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
			reflect.TypeOf((*ax.MessageContext)(nil)).Elem(),
		},
		nil, // no outputs
		aggregateApplySignature,
		aggregateApplySignatureWithMessageContext,
	)
	if err != nil {
		panic(err)
	}

	return a
}

// PersistenceKey returns a unique identifier for the saga.
//
// The persistence key is used to relate persisted data with the saga
// implementation that owns it. Persistence keys should not be changed once
// a saga has active instances.
func (a *Aggregate) PersistenceKey() string {
	return proto.MessageName(a.Prototype)
}

// MessageTypes returns the set of messages that are routed to this saga.
//
// tr is the set of "trigger" messages. If they can not be routed to an
// existing saga instance a new instance is created.
//
// mt is the set of messages that are only routed to existing instances. If
// they can not be routed to an existing instance, the HandleNotFound()
// method is called instead.
func (a *Aggregate) MessageTypes() (tr ax.MessageTypeSet, mt ax.MessageTypeSet) {
	tr = a.Triggers
	return
}

// NewData returns a pointer to a new zero-value instance of the
// saga's data type.
func (a *Aggregate) NewData() Data {
	return proto.Clone(a.Prototype).(Data)
}

// HandleMessage handles a message for a particular saga instance.
func (a *Aggregate) HandleMessage(
	ctx context.Context,
	s ax.Sender,
	mctx ax.MessageContext,
	i Instance,
) (err error) {
	// recordError is a container for errors produced while attempting to record an
	// event.
	type recordError struct{ err error }

	// recover from errors that occur when attempting to record an event
	// re-panic for any other error
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.(recordError); ok {
				err = v.err
			} else {
				panic(r)
			}
		}
	}()

	// wrap any error that occurs while publish in recordError
	rec := func(m ax.Event) {
		if _, err := s.PublishEvent(ctx, m); err != nil {
			panic(recordError{err})
		}
	}

	a.Handle.Dispatch(
		i.Data,
		mctx.Envelope.Message.(ax.Command),
		mctx.Envelope,
		rec,
	)

	return
}

// ApplyEvent updates d to reflect the fact that an event has occurred.
//
// It may panic if env.Message does not implement ax.Event.
func (a *Aggregate) ApplyEvent(d Data, env ax.Envelope) {
	a.Apply.Dispatch(
		d,
		env.Message.(ax.Event),
		env,
	)
}

var (
	aggregateHandleSignature = &typeswitch.Signature{
		Prefix: "Do",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Command)(nil)).Elem(),
			reflect.TypeOf((*ax.EventRecorder)(nil)).Elem(),
		},
	}

	aggregateHandleSignatureWithMessageContext = &typeswitch.Signature{
		Prefix: "Do",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Command)(nil)).Elem(),
			reflect.TypeOf((*ax.MessageContext)(nil)).Elem(),
			reflect.TypeOf((*ax.EventRecorder)(nil)).Elem(),
		},
	}

	aggregateApplySignature = &typeswitch.Signature{
		Prefix: "When",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
		},
	}

	aggregateApplySignatureWithMessageContext = &typeswitch.Signature{
		Prefix: "When",
		In: []reflect.Type{
			reflect.TypeOf((*Data)(nil)).Elem(),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
			reflect.TypeOf((*ax.MessageContext)(nil)).Elem(),
		},
	}
)
