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
	LoadSagaSnapshot(
		ctx context.Context,
		tx persistence.Tx,
		id saga.InstanceID,
	) (i saga.Instance, ok bool, err error)

	// SaveSagaSnapshot saves a snapshot to the store.
	SaveSagaSnapshot(
		ctx context.Context,
		tx persistence.Tx,
		i saga.Instance,
	) error
}
