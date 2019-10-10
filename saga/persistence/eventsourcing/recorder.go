package eventsourcing

import (
	"context"

	"github.com/jmalloc/ax"
)

// Recorder is an implementation of ax.Sender that records published events in
// memory.
type Recorder struct {
	Events []ax.Envelope
	Next   ax.Sender
}

// ExecuteCommand sends a command message.
//
// If ctx contains a message envelope, m is sent as a child of the message in
// that envelope.
func (s *Recorder) ExecuteCommand(
	ctx context.Context,
	m ax.Command,
	opts ...ax.ExecuteOption,
) (ax.Envelope, error) {
	return s.Next.ExecuteCommand(ctx, m, opts...)
}

// PublishEvent sends an event message.
//
// If ctx contains a message envelope, m is sent as a child of the message in
// that envelope.
func (s *Recorder) PublishEvent(
	ctx context.Context,
	m ax.Event,
	opts ...ax.PublishOption,
) (ax.Envelope, error) {
	env, err := s.Next.PublishEvent(ctx, m, opts...)
	if err != nil {
		return ax.Envelope{}, err
	}

	s.Events = append(s.Events, env)

	return env, nil
}
