package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax/persistence"
)

// Repository is an interface for loading and saving saga instances to/from a
// persistent data store.
type Repository interface {
	// LoadSagaInstance fetches a saga instance that has a specific mapping key
	// in its key set.
	//
	// sn is the saga name. k is the message mapping key.
	//
	// If a saga instance is found; ok is true, otherwise it is false. A
	// non-nil error indicates a problem with the store itself.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	LoadSagaInstance(
		ctx context.Context,
		tx persistence.Tx,
		sn, k string,
	) (i Instance, ok bool, err error)

	// SaveSagaInstance persists a saga instance and its associated mapping
	// table to the store as part of tx.
	//
	// It returns an error if the saga instance has been modified since it was
	// loaded, or if there is a problem communicating with the store itself.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	SaveSagaInstance(
		ctx context.Context,
		tx persistence.Tx,
		sn string,
		i Instance,
		ks KeySet,
	) error
}
