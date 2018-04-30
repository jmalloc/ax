package bus

// MessageSender is an interface for sending an outbound message.
import "context"

// MessageSender is an interface for sending messages.
//
// It is implemented by Transport and InboundPipeline.
type MessageSender interface {
	// SendMessage sends a message.
	SendMessage(ctx context.Context, m OutboundEnvelope) error
}

// MessageBuffer is a MessageSender that keeps a collection of sent messages in
// memory.
type MessageBuffer struct {
	Messages []OutboundEnvelope
}

// SendMessage adds m to b.Messages.
func (b *MessageBuffer) SendMessage(ctx context.Context, m OutboundEnvelope) error {
	b.Messages = append(b.Messages, m)
	return nil
}
