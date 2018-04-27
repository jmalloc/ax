package bus

// MessageSender is an interface for sending an outbound message.
import "context"

// MessageSender is an interface for sending messages.
//
// It is implemented by Transport and InboundPipeline.
type MessageSender interface {
	// SendMessage sends a message.
	SendMessage(context.Context, OutboundEnvelope) error
}

// MessageBuffer is a Sender that keeps a collection of sent messages in memory.
type MessageBuffer struct {
	Messages []OutboundEnvelope
}

// DispatchMessage adds m to c.Messages.
func (b *MessageBuffer) DispatchMessage(ctx context.Context, m OutboundEnvelope) error {
	b.Messages = append(b.Messages, m)
	return nil
}
