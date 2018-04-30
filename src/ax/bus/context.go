package bus

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// MessageContext an implementation of ax.MessageContext that sends messages
// using a Sender.
type MessageContext struct {
	context.Context

	Envelope ax.Envelope
	Sender   MessageSender
}

// MessageEnvelope returns the envelope containing the message being handled.
func (c *MessageContext) MessageEnvelope() ax.Envelope {
	return c.Envelope
}

// ExecuteCommand enqueues a command to be executed.
func (c *MessageContext) ExecuteCommand(m ax.Command) error {
	return c.Sender.SendMessage(c, OutboundEnvelope{
		Operation: OpSendUnicast,
		Envelope:  c.Envelope.NewChild(m),
	})
}

// PublishEvent enqueues events to be published.
func (c *MessageContext) PublishEvent(m ax.Event) error {
	return c.Sender.SendMessage(c, OutboundEnvelope{
		Operation: OpSendMulticast,
		Envelope:  c.Envelope.NewChild(m),
	})
}
