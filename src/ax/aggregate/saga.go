package aggregate

import (
	"context"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/saga"
	"github.com/jmalloc/ax/src/internal/visitor"
)

// Aggregate is an alias for saga.Data.
//
// Aggregates are implemented by adding "message handler methods" to a saga.Data
// implementation.
type Aggregate = saga.Data

// Recorder is a function that records the events produced by an aggregate.
type Recorder func(ax.Event)

// New returns a new saga instance that handles messages for the given
// aggregate.
func New(agg Aggregate, opts ...Option) saga.EventedSaga {
	sg := &Saga{
		Prototype: agg,
	}

	for _, opt := range opts {
		opt(sg)
	}

	commandTypes := visitor.MakeAcceptor(
		&sg.CommandHandler,
		reflect.TypeOf((*ax.Command)(nil)).Elem(),
		reflect.TypeOf(agg),
	)

	visitor.MakeAcceptor(
		&sg.EventApplier,
		reflect.TypeOf((*ax.Event)(nil)).Elem(),
		reflect.TypeOf(agg),
	)

	for _, t := range commandTypes {
		sg.CommandTypes = sg.CommandTypes.Add(
			ax.TypeOf(
				reflect.Zero(t).Interface().(ax.Message),
			),
		)
	}

	return sg
}

// Saga is an implementation of saga.Saga that wraps an AggregateRoot.
type Saga struct {
	saga.MapByInstanceID
	saga.ErrorIfNotFound

	Prototype      Aggregate
	Identifier     Identifier
	CommandTypes   ax.MessageTypeSet
	CommandHandler func(Aggregate, ax.Command, Recorder)
	EventApplier   func(Aggregate, ax.Event)
}

// SagaName returns a unique name for the saga.
//
// The saga name is used to relate saga instances to the saga implementation
// that manages them. For that reason, saga names should not be changed when
// there are active saga instances.
func (sg *Saga) SagaName() string {
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
	return sg.CommandTypes, ax.MessageTypeSet{}
}

// GenerateInstanceID returns the saga ID to use for a new instance.
//
// It is called when a "trigger" message is received and there is no
// existing saga instance. env contains the "trigger" message.
func (sg *Saga) GenerateInstanceID(ctx context.Context, env ax.Envelope) (id saga.InstanceID, err error) {
	m := env.Message.(ax.Command)
	v, err := sg.Identifier.AggregateID(m)
	if err != nil {
		return
	}

	err = id.Parse(v)
	return
}

// NewData returns a pointer to a new zero-value instance of the
// saga's data type.
func (sg *Saga) NewData() saga.Data {
	return proto.Clone(sg.Prototype).(saga.Data)
}

// MappingKeyForMessage returns the key used to locate the saga instance
// to which the given message is routed, if any.
//
// If ok is false the message is ignored; otherwise, the message is routed
// to the saga instance that contains k in its associated key set.
//
// New saga instances are created when no matching instance can be found
// and the message is declared as a "trigger" by the saga's MessageTypes()
// method; otherwise, HandleNotFound() is called.
func (sg *Saga) MappingKeyForMessage(ctx context.Context, env ax.Envelope) (k string, ok bool, err error) {
	m := env.Message.(ax.Command)
	ok = true
	k, err = sg.Identifier.AggregateID(m)
	return
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

	sg.CommandHandler(
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
	sg.EventApplier(d, m)
}
