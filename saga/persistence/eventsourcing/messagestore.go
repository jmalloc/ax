package eventsourcing

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/messagestore"
	"github.com/jmalloc/ax/persistence"
	"github.com/jmalloc/ax/saga"
)

// streamName returns the message store stream name that contains the events for
// the given instance.
func streamName(id saga.InstanceID) string {
	return "saga:" + id.Get()
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
		ok, err := s.TryNext(ctx)
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
