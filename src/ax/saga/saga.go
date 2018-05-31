package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// A Saga is a stateful message handler.
//
// They are typically used to model "long-running" business processes. They are
// the foundation on which aggregates and workflows are built.
//
// Each saga can have multiple instances, represented by the saga.Instance
// struct. Each instance has associated application-defined data, represented by
// the saga.Data interface.
//
// For each saga, an inbound message is always routed to one saga instance.
//
// Saga instances are persisted using an implementation of the Repository
// interface, which is typically provided by a specific persistence
// implementation.
type Saga interface {
	// SagaName returns a unique name for the saga.
	//
	// The saga name is used to relate saga instances to the saga implementation
	// that manages them. For that reason, saga names should not be changed when
	// there are active saga instances.
	SagaName() string

	// MessageTypes returns the set of messages that are routed to this saga.
	//
	// tr is the set of "trigger" messages. If they can not be routed to an
	// existing saga instance a new instance is created.
	//
	// mt is the set of messages that are only routed to existing instances. If
	// they can not be routed to an existing instance, the HandleNotFound()
	// method is called instead.
	MessageTypes() (tr ax.MessageTypeSet, mt ax.MessageTypeSet)

	// GenerateInstanceID returns the saga ID to use for a new instance.
	//
	// It is called when a "trigger" message is received and there is no
	// existing saga instance. env contains the "trigger" message.
	GenerateInstanceID(ctx context.Context, env ax.Envelope) (id InstanceID, err error)

	// NewData returns a pointer to a new zero-value instance of the
	// saga's data type.
	NewData() Data

	// MappingKeyForMessage returns the key used to locate the saga instance
	// to which the given message is routed, if any.
	//
	// If ok is false the message is ignored; otherwise, the message is routed
	// to the saga instance that contains k in its associated key set.
	//
	// New saga instances are created when no matching instance can be found
	// and the message is declared as a "trigger" by the saga's MessageTypes()
	// method; otherwise, HandleNotFound() is called.
	MappingKeyForMessage(ctx context.Context, env ax.Envelope) (k string, ok bool, err error)

	// MappingKeysForInstance returns the set of mapping keys associated with
	// the given instance.
	//
	// When a message is received, a mapping key is produced by calling
	// MappingKeyForMessage(). The message is routed to the saga instance that
	// contains this key in its key set.
	//
	// Key sets must be disjoint. That is, no two instances of the same saga
	// may share any keys.
	//
	// The key set is rebuilt whenever a message is received. It is persisted
	// alongside the saga instance by the Repository.
	//
	// Extra care must be taken when introducing a new key to the set, as the key
	// set for existing saga instances will not be updated until they next receive
	// a message.
	MappingKeysForInstance(context.Context, Instance) ([]string, error)

	// HandleMessage handles a message for a particular saga instance.
	HandleMessage(context.Context, ax.Sender, ax.Envelope, Instance) error

	// HandleNotFound handles a message that is intended for a saga instance
	// that could not be found.
	HandleNotFound(context.Context, ax.Sender, ax.Envelope) error
}

// EventedSaga is a saga that only mutates its data when an event occurs.
//
// CRUD sagas may be evented or non-evented, but eventsourced sagas are always
// evented.
//
// Implementors should take care not to mutate the saga data directly inside the
// saga HandleMessage() method, only in ApplyEvent().
type EventedSaga interface {
	Saga

	// ApplyEvent updates d to reflect the fact that an event has occurred.
	//
	// It may panic if env.Message does not implement ax.Event.
	ApplyEvent(d Data, env ax.Envelope)
}
