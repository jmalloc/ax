package eventsourcing

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/messagestore"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// streamName returns the message store stream name that contains the events for
// the given instance.
func streamName(id saga.InstanceID) string {
	return "saga:" + id.Get()
}

// appendEvents appends the events in envs to the message stream for the given
// saga instance.
//
// It returns an error if i.Revision is not the next free offset in the stream.
// It panics if envs contains messages that do not implement ax.Event.
func appendEvents(
	ctx context.Context,
	tx persistence.Tx,
	ms messagestore.Store,
	i saga.Instance,
	envs []ax.Envelope,
) error {
	for _, env := range envs {
		_ = env.Message.(ax.Event) // panic if not an event
	}

	return ms.AppendMessages(
		ctx,
		tx,
		streamName(i.InstanceID),
		uint64(i.Revision),
		envs,
	)
}

// applyEvents calls Data.ApplyEvent for each event in a saga's message stream.
func applyEvents(
	ctx context.Context,
	tx persistence.Tx,
	ms messagestore.Store,
	sg saga.EventedSaga,
	i *saga.Instance,
) error {
	s, ok, err := ms.OpenStream(
		ctx,
		tx.DataStore(),
		streamName(i.InstanceID),
		uint64(i.Revision),
	)
	if !ok || err != nil {
		return err
	}

	for {
		ok, err := s.Next(ctx)
		if !ok || err != nil {
			return err
		}

		env, err := s.Get(ctx)
		if err != nil {
			return err
		}

		if _, ok := env.Message.(ax.Event); !ok {
			return fmt.Errorf(
				"event stream for saga instance %s contains non-event message %s",
				i.InstanceID.Get(),
				env.MessageID.Get(),
			)
		}

		sg.ApplyEvent(i.Data, env)
		i.Revision++
	}
}
