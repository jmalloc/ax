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
	// It returns an error if the instance does not exist, or a problem occurs
	// with the underlying data store.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	LoadSagaInstance(ctx context.Context, tx persistence.Tx, id saga.InstanceID) (saga.Instance, error)

	// SaveSagaInstance persists a saga instance.
	//
	// It returns an error if the saga instance has been modified since it was
	// loaded, or a problem occurs with the underlying data store.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	SaveSagaInstance(ctx context.Context, tx persistence.Tx, i saga.Instance) error
}
