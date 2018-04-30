package bus

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// MessageContext is an implementation of ax.MessageContext that sends messages
// using a MessageSender.
type MessageContext struct {
	context.Context

	Envelope ax.Envelope
	Sender   MessageSender
}

// MessageEnvelope returns the envelope containing the message being handled.
func (c *MessageContext) MessageEnvelope() ax.Envelope {
	return c.Envelope
}

// ExecuteCommand sends a command message.
//
// Commands are routed to a single endpoint as per the routing rules of the
// outbound message pipeline.
func (c *MessageContext) ExecuteCommand(m ax.Command) error {
	return c.Sender.SendMessage(c, OutboundEnvelope{
		Operation: OpSendUnicast,
		Envelope:  c.Envelope.NewChild(m),
	})
}

// PublishEvent sends an event message.
//
// Events are routed to endpoints that subscribe to messages of that type.
func (c *MessageContext) PublishEvent(m ax.Event) error {
	return c.Sender.SendMessage(c, OutboundEnvelope{
		Operation: OpSendMulticast,
		Envelope:  c.Envelope.NewChild(m),
	})
}
