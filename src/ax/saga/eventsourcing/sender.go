package eventsourcing

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Sender is an implementation of ax.Sender that applies published
// events to the state of an eventsourced saga.
type Sender struct {
	Data   Data
	Events []ax.Event
	Next   ax.Sender
}

// ExecuteCommand sends a command message.
//
// If ctx contains a message envelope, m is sent as a child of the message in
// that envelope.
func (s *Sender) ExecuteCommand(ctx context.Context, m ax.Command) error {
	return s.Next.ExecuteCommand(ctx, m)
}

// PublishEvent sends an event message.
//
// If ctx contains a message envelope, m is sent as a child of the message in
// that envelope.
func (s *Sender) PublishEvent(ctx context.Context, m ax.Event) error {
	s.Events = append(s.Events, m)
	s.Data.ApplyEvent(m)
	return s.Next.PublishEvent(ctx, m)
}
