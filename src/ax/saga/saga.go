package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Saga is an interface for handling messages associated with a particular saga
// instance.
//
// A Saga is essentially a stateful message handler where persistence of the
// saga state is managed by the framework. The state is represented by the
// Instance interface. Each saga may produce many instances.
type Saga interface {
	// SagaName returns a unique name for the saga.
	// The saga name is used to locate instances of the saga. It must not be
	// changed while there are active instances.
	SagaName() string

	// MessageTypes returns the set of messages that are routed to this saga.
	//
	// tr is the set of "trigger" messages that will cause a new instance to be
	// created. mt is the set of messages that are only routed to existing
	// instances (or the not-found handler).
	MessageTypes() (tr ax.MessageTypeSet, mt ax.MessageTypeSet)

	// NewInstance returns a new saga instance.
	NewInstance(context.Context, ax.Envelope) (InstanceID, Data, error)

	// BuildMappingTable returns the message mapping table to use for the given
	// saga instance.
	//
	// Mapping tables are used to correlate incoming messages with the saga
	// instance they are routed to.
	//
	// The mapping table is rebuilt each time an instance receives a message. Care
	// should be taken when adding new keys to the mapping table, as the persisted
	// mapping tables for existing instances will not include that key until they
	// next receive a message.
	BuildMappingTable(context.Context, Instance) (map[string]string, error)

	// MapMessage returns the key and value to use to locate the saga instance
	// for the given message.
	MapMessage(context.Context, ax.Envelope) (string, string, error)

	// HandleMessage handles a message for a particular saga instance.
	HandleMessage(context.Context, ax.Sender, ax.Envelope, Instance) error

	// HandleNotFound handles a message that is intended for a saga instance
	// that could not be found.
	HandleNotFound(context.Context, ax.Sender, ax.Envelope) error
}
