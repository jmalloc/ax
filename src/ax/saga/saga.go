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
	// PersistenceKey returns a unique identifier for the saga.
	//
	// The persistence key is used to relate persisted data with the saga
	// implementation that owns it. Persistence keys should not be changed once
	// a saga has active instances.
	PersistenceKey() string

	// MessageTypes returns the set of messages that are routed to this saga.
	//
	// tr is the set of "trigger" messages. If they can not be routed to an
	// existing saga instance a new instance is created.
	//
	// mt is the set of messages that are only routed to existing instances. If
	// they can not be routed to an existing instance, the HandleNotFound()
	// method is called instead.
	MessageTypes() (tr ax.MessageTypeSet, mt ax.MessageTypeSet)

	// NewData returns a pointer to a new zero-value instance of the
	// saga's data type.
	NewData() Data

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
