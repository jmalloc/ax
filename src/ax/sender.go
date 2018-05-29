package ax

import (
	"context"
)

// Sender is an interface for sending messages.
type Sender interface {
	// ExecuteCommand sends a command message.
	//
	// Commands are routed to a single endpoint as per the routing rules of the
	// outbound message pipeline.
	ExecuteCommand(context.Context, Command) (Envelope, error)

	// PublishEvent sends an event message.
	//
	// Events are routed to endpoints that subscribe to messages of that type.
	PublishEvent(context.Context, Event) (Envelope, error)
}
