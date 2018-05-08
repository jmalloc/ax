package bus

import "context"

// MessageSource is an interface that produces inbound message envelopes.
type MessageSource interface {
	// Produce returns the next message from the source.
	// It blocks until a message is available, or ctx is canceled.
	Produce(ctx context.Context) (InboundEnvelope, error)
}
