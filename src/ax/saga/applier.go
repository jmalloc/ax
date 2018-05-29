package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Applier is an implementation of ax.Sender that applies published
// events to saga data for evented sagas.
type Applier struct {
	Saga EventedSaga
	Data Data

	Next ax.Sender
}

// ExecuteCommand sends a command message.
//
// If ctx contains a message envelope, m is sent as a child of the message in
// that envelope.
func (s *Applier) ExecuteCommand(ctx context.Context, m ax.Command) (ax.Envelope, error) {
	return s.Next.ExecuteCommand(ctx, m)
}

// PublishEvent sends an event message.
//
// If ctx contains a message envelope, m is sent as a child of the message in
// that envelope.
func (s *Applier) PublishEvent(ctx context.Context, m ax.Event) (ax.Envelope, error) {
	env, err := s.Next.PublishEvent(ctx, m)
	if err != nil {
		return ax.Envelope{}, err
	}

	s.Saga.ApplyEvent(s.Data, env)

	return env, nil
}
