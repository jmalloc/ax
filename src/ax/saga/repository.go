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
	// mt is the the type of the message being processed, and k is the mapping
	// key for that message that identifies the saga to load.
	//
	// ok is true if the instance is found, err is non-nil if there is a problem
	// communicating with the store itself.
	LoadSagaInstance(
		ctx context.Context,
		tx persistence.Tx,
		req LoadRequest,
	) (res LoadResult, ok bool, err error)

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
		req SaveRequest,
	) error
}

// LoadRequest contains information used to locate the saga instance that
// is associated with a specific message.
type LoadRequest struct {
	SagaName    string
	MessageType ax.MessageType
	MappingKey  string
}

// LoadResult contains the result of loading a saga.
type LoadResult struct {
	InstanceID      InstanceID
	CurrentRevision uint64
	Instance        Instance
}

// SaveRequest contains information used to save a saga instance.
type SaveRequest struct {
	SagaName        string
	InstanceID      InstanceID
	CurrentRevision uint64
	Instance        Instance
	MappingTable    MappingTable
}
