package eventsourcing

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// EventStore is an interface for reading and writing streams of events.
type EventStore interface {
	AppendEvents(
		ctx context.Context,
		tx persistence.Tx,
		id saga.InstanceID,
		rev saga.Revision,
		ev []ax.Event,
	) error

	OpenStream(
		ctx context.Context,
		tx persistence.Tx,
		id saga.InstanceID,
		rev saga.Revision,
	) (EventStream, error)
}

// EventStream is a stream of events produced by an eventsourced saga.
type EventStream interface {
	// Next advances the stream to the next event.
	// It returns false if there are no more events.
	Next(ctx context.Context) (bool, error)

	// Get returns the event at the current location in the stream.
	Get(ctx context.Context) (ax.Event, error)

	// Close closes the stream.
	Close() error
}

func applyEvents(
	ctx context.Context,
	tx persistence.Tx,
	es EventStore,
	i *saga.Instance,
) error {
	s, err := es.OpenStream(
		ctx,
		tx,
		i.InstanceID,
		i.Revision+1,
	)
	if err != nil {
		return err
	}

	data := i.Data.(Data)

	for {
		ok, err := s.Next(ctx)
		if !ok || err != nil {
			return err
		}

		ev, err := s.Get(ctx)
		if err != nil {
			return err
		}

		data.ApplyEvent(ev)
		i.Revision++
	}
}
