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
	// MessageTypes returns the set of messages that are routed to this saga.
	//
	// tr is the set of "trigger" messages that will cause a new instance to be
	// created. mt is the set of messages that are only routed to existing
	// instances (or the not-found handler).
	MessageTypes() (tr ax.MessageTypeSet, mt ax.MessageTypeSet)

	// MapMessage returns a mapping key for the given message.
	MapMessage(ax.Message) MappingKey

	// MapInstance returns a mapping key to use for the given message
	// type and saga instance.
	MapInstance(ax.MessageType, Instance) MappingKey

	// InitialState returns a new saga value.
	InitialState() Instance

	// HandleMessage handles a message for a particular saga instance.
	HandleMessage(context.Context, ax.Sender, ax.Envelope, Instance) error

	// HandleNotFound handles a message that is intended for a saga instance
	// that could not be found.
	HandleNotFound(context.Context, ax.Sender, ax.Envelope) error
}
