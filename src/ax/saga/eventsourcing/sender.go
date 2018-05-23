package eventsourcing

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Sender is an implementation of ax.Sender that applies published
// events to the state of an eventsourced saga.
type Sender struct {
	Data   Data
	Events []ax.Envelope
	Next   ax.Sender
}

// ExecuteCommand sends a command message.
//
// If ctx contains a message envelope, m is sent as a child of the message in
// that envelope.
func (s *Sender) ExecuteCommand(ctx context.Context, m ax.Command) (ax.Envelope, error) {
	return s.Next.ExecuteCommand(ctx, m)
}

// PublishEvent sends an event message.
//
// If ctx contains a message envelope, m is sent as a child of the message in
// that envelope.
func (s *Sender) PublishEvent(ctx context.Context, m ax.Event) (ax.Envelope, error) {
	env, err := s.Next.PublishEvent(ctx, m)
	if err != nil {
		return ax.Envelope{}, err
	}

	s.Events = append(s.Events, env)
	s.Data.ApplyEvent(env)

	return env, nil
}
