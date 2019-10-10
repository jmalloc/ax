package crud

import (
	"context"

	"github.com/jmalloc/ax/persistence"
	"github.com/jmalloc/ax/saga"
)

// Repository is an interface for loading and saving saga instances.
type Repository interface {
	// LoadSagaInstance fetches a saga instance by its ID.
	//
	// It returns an false if the instance does not exist. It returns an error
	// if a problem occurs with the underlying data store.
	//
	// It returns an error if the instance is found, but belongs to a different
	// saga, as identified by pk, the saga's persistence key.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	LoadSagaInstance(
		ctx context.Context,
		tx persistence.Tx,
		pk string,
		id saga.InstanceID,
	) (saga.Instance, bool, error)

	// SaveSagaInstance persists a saga instance.
	//
	// It returns an error if i.Revision is not the current revision of the
	// instance as it exists within the store, or a problem occurs with the
	// underlying data store.
	//
	// It returns an error if the instance belongs to a different saga, as
	// identified by pk, the saga's persistence key.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	SaveSagaInstance(
		ctx context.Context,
		tx persistence.Tx,
		pk string,
		i saga.Instance,
	) error

	// DeleteSagaInstance deletes a saga instance.
	//
	// It returns an error if i.Revision is not the current revision of the
	// instance as it exists within the store, or a problem occurs with the
	// underlying data store.
	//
	// It returns an error if the instance belongs to a different saga, as
	// identified by pk, the saga's persistence key.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	DeleteSagaInstance(
		ctx context.Context,
		tx persistence.Tx,
		pk string,
		i saga.Instance,
	) error
}
