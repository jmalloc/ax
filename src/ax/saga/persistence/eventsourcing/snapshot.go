package eventsourcing

import (
	"context"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// SnapshotRepository is an interface for loading and saving snapshots of
// eventsourced saga data.
type SnapshotRepository interface {
	// LoadSagaSnapshot loads the latest available snapshot from the store.
	//
	// It returns an error if a snapshot of this instance is found, but belongs to
	// a different saga, as identified by pk, the saga's persistence key.
	LoadSagaSnapshot(
		ctx context.Context,
		tx persistence.Tx,
		pk string,
		id saga.InstanceID,
	) (i saga.Instance, ok bool, err error)

	// SaveSagaSnapshot saves a snapshot to the store.
	//
	// The implementation may return an error if a snapshot for this instance
	// already exists, but belongs to a different saga, as identified by pk, the
	// saga's persistence key.
	SaveSagaSnapshot(
		ctx context.Context,
		tx persistence.Tx,
		pk string,
		i saga.Instance,
	) error

	// DeleteSagaSnapshots deletes any snapshots associated with a saga instance.
	//
	// The implementation may return an error if snapshots for this instance
	// already exists, but belong to a different saga, as identified by pk, the
	// saga's persistence key.
	DeleteSagaSnapshots(
		ctx context.Context,
		tx persistence.Tx,
		pk string,
		id saga.InstanceID,
	) error
}
