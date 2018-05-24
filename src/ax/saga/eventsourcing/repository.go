package eventsourcing

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// InstanceRepository is an interface for loading and saving eventsourced saga
// instances.
type InstanceRepository interface {
	// LoadSagaInstance fetches a saga instance by its ID.
	//
	// If a saga instance is found; ok is true, otherwise it is false. A
	// non-nil error indicates a problem with the store itself.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	LoadSagaInstance(
		ctx context.Context,
		tx persistence.Tx,
		id saga.InstanceID,
		d saga.EventedData,
	) (saga.Instance, error)

	// SaveSagaInstance persists a saga instance.
	//
	// It returns an error if the saga instance has been modified since it was
	// loaded, or if there is a problem communicating with the store itself.
	//
	// It panics if envs contains any messages that do not implement ax.Event.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	SaveSagaInstance(
		ctx context.Context,
		tx persistence.Tx,
		i saga.Instance,
		envs []ax.Envelope,
	) error
}

// DefaultSnapshotFrequency is the default number of revisions to allow between
// storing snapshots.
const DefaultSnapshotFrequency = 1000

// StandardInstanceRepository is an implementation of InstanceRepository that
// uses a persistence.MessageStore and an optional SnapshotRepository to store
// event-sourced saga instances.
type StandardInstanceRepository struct {
	MessageStore      persistence.MessageStore
	Snapshots         SnapshotRepository
	SnapshotFrequency saga.Revision
}

// LoadSagaInstance fetches a saga instance by its ID.
//
// It first attempts to load a snapshot of the saga data, before fetching
// events from the message store.
//
// If a saga instance is found; ok is true, otherwise it is false. A
// non-nil error indicates a problem with the store itself.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r *StandardInstanceRepository) LoadSagaInstance(
	ctx context.Context,
	tx persistence.Tx,
	id saga.InstanceID,
	d saga.EventedData,
) (saga.Instance, error) {
	i := saga.Instance{
		InstanceID: id,
		Data:       d,
	}

	if r.Snapshots != nil {
		var err error
		i, _, err = r.Snapshots.LoadSagaSnapshot(ctx, tx, id)
		if err != nil {
			return i, err
		}
	}

	return i, applyEvents(ctx, tx, r.MessageStore, &i)
}

// SaveSagaInstance persists a saga instance.
//
// If the save operation causes the saga revision to pass the snapshot frequency
// threshold, a new snapshot is stored.
//
// It returns an error if the saga instance has been modified since it was
// loaded, or if there is a problem communicating with the store itself.
//
// It panics if envs contains any messages that do not implement ax.Event.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r *StandardInstanceRepository) SaveSagaInstance(
	ctx context.Context,
	tx persistence.Tx,
	i saga.Instance,
	envs []ax.Envelope,
) error {
	if err := appendEvents(ctx, tx, r.MessageStore, i, envs); err != nil {
		return err
	}

	prev := i.Revision
	i.Revision += saga.Revision(len(envs))

	if r.shouldSnapshot(prev, i.Revision) {
		return r.Snapshots.SaveSagaSnapshot(ctx, tx, i)
	}

	return nil
}

// shouldSnapshot returns true if a new snapshot should be stored.
func (r *StandardInstanceRepository) shouldSnapshot(before, after saga.Revision) bool {
	if r.Snapshots == nil {
		return false
	}

	freq := r.SnapshotFrequency
	if freq == 0 {
		freq = DefaultSnapshotFrequency
	}

	return (before / freq) != (after / freq)
}

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
	ms persistence.MessageStore,
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
	ms persistence.MessageStore,
	i *saga.Instance,
) error {
	s, err := ms.OpenStream(
		ctx,
		tx,
		streamName(i.InstanceID),
		uint64(i.Revision),
	)
	if err != nil {
		return err
	}

	data := i.Data.(saga.EventedData)

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

		data.ApplyEvent(env)
		i.Revision++
	}
}
