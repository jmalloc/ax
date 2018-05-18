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
// struct. Each instance as associated application-defined data, represented by
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

	// NewInstance returns a new saga instance.
	//
	// It is called when a "trigger" message is received and there is no
	// existing saga instance. env contains the "trigger" message.
	//
	// If err is nil, id must be a valid InstanceID, and d must be non-nil.
	NewInstance(ctx context.Context, env ax.Envelope) (id InstanceID, d Data, err error)

	// MappingKeyForMessage returns the key used to locate the saga instance
	// to which the given message is routed.
	//
	// The message is routed to the saga instance that contains k in its
	// associated key set.
	//
	// If no saga instance is found and the message is a "trigger" message, a
	// new instance is created; otherwise, HandleNotFound() is called.
	MappingKeyForMessage(ctx context.Context, env ax.Envelope) (k string, err error)

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
	// The key set is rebuild whenever a message is received. It is persisted
	// alongside the saga instance by the Repository.
	//
	// Extra care must be taken when introducing a new key to the set, as the key
	// set for existing saga instances will not be updated until they next receive
	// a message.
	MappingKeysForInstance(context.Context, Instance) (KeySet, error)

	// HandleMessage handles a message for a particular saga instance.
	HandleMessage(context.Context, ax.Sender, ax.Envelope, Instance) error

	// HandleNotFound handles a message that is intended for a saga instance
	// that could not be found.
	HandleNotFound(context.Context, ax.Sender, ax.Envelope) error
}
