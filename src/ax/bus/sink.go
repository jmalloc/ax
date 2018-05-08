package bus

import "context"

// MessageSink is an interface that accepts outbound message envelopes as input.
type MessageSink interface {
	// Accept processes the message encapsulated in env.
	Accept(ctx context.Context, env OutboundEnvelope) error
}

// BufferedSink is a MessageSink that buffers message envelopes in memory.
// in memory.
type BufferedSink struct {
	Envelopes []OutboundEnvelope
}

// Accept adds env to b.Envelopes.
func (b *BufferedSink) Accept(ctx context.Context, env OutboundEnvelope) error {
	b.Envelopes = append(b.Envelopes, env)
	return nil
}
