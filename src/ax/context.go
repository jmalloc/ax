package ax

import (
	"context"
)

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

// BindContext returns a new MessageContext that forwards messages to mc,
// but uses ctx for context.Context related operations.
func BindContext(ctx context.Context, mc MessageContext) MessageContext {
	return boundContext{ctx, mc}
}

// boundContext is an implementation of MessageContext that replaces the
// context.Context features of an existing MessageContext with a specific
// context.Context instance.
type boundContext struct {
	context.Context
	parent MessageContext
}

func (c boundContext) MessageEnvelope() Envelope {
	return c.parent.MessageEnvelope()
}

func (c boundContext) ExecuteCommand(m Command) error {
	return c.parent.ExecuteCommand(m)
}

func (c boundContext) PublishEvent(m Event) error {
	return c.parent.PublishEvent(m)
}
