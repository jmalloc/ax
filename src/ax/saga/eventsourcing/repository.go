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

// MessageStoreRepository is an implementation of Repository that loads
// eventsourced saga instances from an event store, without using any snapshots.
type MessageStoreRepository struct {
	MessageStore persistence.MessageStore
}

// LoadSagaInstance fetches a saga instance by its ID.
//
// If a saga instance is found; ok is true, otherwise it is false. A
// non-nil error indicates a problem with the store itself.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r *MessageStoreRepository) LoadSagaInstance(
	ctx context.Context,
	tx persistence.Tx,
	id saga.InstanceID,
	d Data,
) (saga.Instance, error) {
	i := saga.Instance{
		InstanceID: id,
		Data:       d,
	}

	return i, applyEvents(ctx, tx, r.MessageStore, &i)
}

// SaveSagaInstance persists a saga instance.
//
// It returns an error if the saga instance has been modified since it was
// loaded, or if there is a problem communicating with the store itself.
//
// It panics if envs contains any messages that do not implement ax.Event.
//
// It panics if the repository is not able to enlist in tx because it uses a
// different underlying storage system.
func (r *MessageStoreRepository) SaveSagaInstance(
	ctx context.Context,
	tx persistence.Tx,
	i saga.Instance,
	envs []ax.Envelope,
) error {
	return appendEvents(ctx, tx, r.MessageStore, i, envs)
}

// DefaultSnapshotFrequency is the default number of revisions to allow between
// storing snapshots.
const DefaultSnapshotFrequency = 1000

// SnapshottingRepository is an implementation of Repository that uses an
// event store and snaphot repository to store eventsourced saga instances.
type SnapshottingRepository struct {
	MessageStore persistence.MessageStore
	Snapshots    SnapshotRepository
	Frequency    saga.Revision
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
func (r *SnapshottingRepository) SaveSagaInstance(
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
func (r *SnapshottingRepository) shouldSnapshot(before, after saga.Revision) bool {
	freq := r.Frequency
	if freq == 0 {
		freq = DefaultSnapshotFrequency
	}

	return (before / freq) != (after / freq)
}
