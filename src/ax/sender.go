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
	ExecuteCommand(context.Context, Command, ...ExecuteOption) (Envelope, error)

	// PublishEvent sends an event message.
	//
	// Events are routed to endpoints that subscribe to messages of that type.
	PublishEvent(context.Context, Event, ...PublishOption) (Envelope, error)
}

// ExecuteOption is configures an envelope containing a command message to
// exhibit some specific behavior.
type ExecuteOption interface {
	ApplyExecuteOption(env *Envelope) error
}

// PublishOption is configures an envelope containing an event message to
// exhibit some specific behavior.
type PublishOption interface {
	ApplyPublishOption(env *Envelope) error
}

// SendOption is an option that can be used for both commands and events.
type SendOption interface {
	ExecuteOption
	PublishOption
}

// EventRecorder is a function that records the occurrence of events.
type EventRecorder func(Event)

// CommandExecutor is a function that queues a command to be executed.
type CommandExecutor func(Command, ...ExecuteOption)
