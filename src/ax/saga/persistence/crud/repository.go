package crud

import (
	"context"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// Repository is an interface for loading and saving saga instances.
type Repository interface {
	// LoadSagaInstance fetches a saga instance by its ID.
	//
	// It returns an false if the instance does not exist. It returns an error
	// if a problem occurs with the underlying data store.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	LoadSagaInstance(
		ctx context.Context,
		tx persistence.Tx,
		id saga.InstanceID,
	) (saga.Instance, bool, error)

	// SaveSagaInstance persists a saga instance.
	//
	// It returns an error if i.Revision is not the current revision of the
	// instance as it exists within the store, or a problem occurs with the
	// underlying data store.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	SaveSagaInstance(ctx context.Context, tx persistence.Tx, i saga.Instance) error
}
