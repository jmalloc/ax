package eventsourcing

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// Repository is an interface for loading and saving eventsourced saga instances.
type Repository interface {
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
		d Data,
	) (saga.Instance, error)

	SaveSagaInstance(
		ctx context.Context,
		tx persistence.Tx,
		i saga.Instance,
		ev []ax.Event,
	) error
}

// EventStoreRepository is an implementation of Repository that loads
// eventsourced saga instances from an event store, without using any snapshots.
type EventStoreRepository struct {
	EventStore EventStore
}

func (r *EventStoreRepository) LoadSagaInstance(
	ctx context.Context,
	tx persistence.Tx,
	id saga.InstanceID,
	d Data,
) (saga.Instance, error) {
	i := saga.Instance{
		InstanceID: id,
		Data:       d,
	}

	return i, applyEvents(ctx, tx, r.EventStore, &i)
}

func (r *EventStoreRepository) SaveSagaInstance(
	ctx context.Context,
	tx persistence.Tx,
	i saga.Instance,
	ev []ax.Event,
) error {
	return r.EventStore.AppendEvents(ctx, tx, i.InstanceID, i.Revision, ev)
}

const DefaultSnapshotFrequency = 1000

// SnapshottingRepository is an implementation of Repository that uses an
// event store and snaphot repository to store eventsourced saga instances.
type SnapshottingRepository struct {
	EventStore        EventStore
	Snapshots         SnapshotRepository
	SnapshotFrequency saga.Revision
}

func (r *SnapshottingRepository) LoadSagaInstance(
	ctx context.Context,
	tx persistence.Tx,
	id saga.InstanceID,
	d Data,
) (saga.Instance, error) {
	i, ok, err := r.Snapshots.LoadSagaSnapshot(ctx, tx, id)
	if err != nil {
		return i, err
	}

	if !ok {
		i = saga.Instance{
			InstanceID: id,
			Data:       d,
		}
	}

	return i, applyEvents(ctx, tx, r.EventStore, &i)
}

func (r *SnapshottingRepository) SaveSagaInstance(
	ctx context.Context,
	tx persistence.Tx,
	i saga.Instance,
	ev []ax.Event,
) error {
	if err := r.EventStore.AppendEvents(ctx, tx, i.InstanceID, i.Revision, ev); err != nil {
		return err
	}

	prev := i.Revision
	i.Revision += saga.Revision(len(ev))

	if r.shouldSnapshot(prev, i.Revision) {
		return r.Snapshots.SaveSagaSnapshot(ctx, tx, i)
	}

	return nil
}

func (r *SnapshottingRepository) shouldSnapshot(before, after saga.Revision) bool {
	freq := r.SnapshotFrequency
	if freq == 0 {
		freq = DefaultSnapshotFrequency
	}

	return (before / freq) != (after / freq)
}
