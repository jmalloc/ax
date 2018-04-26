package ax

import "context"

// MessageContext is a specialization of context.Context used to send messages
// within the context of handling a message.
//
// Messages sent via a context are configured as "children" of the message being
// handled in that context.
type MessageContext interface {
	context.Context

	// MessageEnvelope returns the envelope containing the message being handled.
	MessageEnvelope() Envelope

	// ExecuteCommand sends a command message.
	//
	// Commands are routed to a single endpoint as per the routing rules of the
	// outbound message pipeline.
	ExecuteCommand(Command) error

	// PublishEvent sends an event message.
	//
	// Events are routed to endpoints that subscribe to messages of that type.
	PublishEvent(Event) error
}
