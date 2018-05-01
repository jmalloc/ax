package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// Repository is an interface for loading and saving saga instances to/from a
// persistent data store.
type Repository interface {
	// LoadSagaInstance fetches a saga instance from the store based on a
	// mapping key for a particular message type.
	//
	// ok is true if the instance is found, in which case i is populated with
	// data from the store.
	//
	// err is non-nil if there is a problem communicating with the store itself.
	LoadSagaInstance(
		ctx context.Context,
		mt ax.MessageType,
		k MappingKey,
		i Instance,
	) (ok bool, err error)

	// SaveSagaInstance persists a saga instance and its associated mapping
	// table to the store as part of tx.
	//
	// It returns an error if the saga instance has been modified since it was
	// loaded, or if there is a problem communicating with the store itself.
	//
	// Save() panics if the repository is not able to enlist in tx because it
	// uses a different underlying storage system.
	SaveSagaInstance(
		ctx context.Context,
		tx persistence.Tx,
		i Instance,
		t MappingTable,
	) error
}
