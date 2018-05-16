package endpoint

import (
	"context"
	"sync"
)

// MessageSink is an interface that accepts outbound message envelopes as input.
type MessageSink interface {
	// Accept processes the message encapsulated in env.
	Accept(ctx context.Context, env OutboundEnvelope) error
}

// BufferedSink is a MessageSink that buffers message envelopes in memory.
type BufferedSink struct {
	m   sync.RWMutex
	env []OutboundEnvelope
}

// Accept buffers env in memory.
func (s *BufferedSink) Accept(ctx context.Context, env OutboundEnvelope) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.env = append(s.env, env)

	return nil
}

// Reset removes the buffered message envelopes.
func (s *BufferedSink) Reset() {
	s.m.Lock()
	defer s.m.Unlock()

	s.env = nil
}

// Envelopes returns the message envelopes that have been buffered.
func (s *BufferedSink) Envelopes() []OutboundEnvelope {
	s.m.RLock()
	defer s.m.RUnlock()

	return append([]OutboundEnvelope(nil), s.env...)
}

// TakeEnvelopes returns the message envelopes that have been buffered and
// resets the sink in a single operation.
func (s *BufferedSink) TakeEnvelopes() []OutboundEnvelope {
	s.m.Lock()
	defer s.m.Unlock()

	env := s.env
	s.env = nil

	return env
}
