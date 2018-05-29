package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax/persistence"
)

// Mapper is an interface for finding saga instances by their mapping key.
type Mapper interface {
	// FindByKey returns the instance ID of the saga instance that handles
	// messages with a specific mapping key.
	//
	// sn is the name of the saga, and k is the message's mapping key.
	//
	// ok is false if no saga instance is found.
	FindByKey(
		ctx context.Context,
		tx persistence.Tx,
		sn, k string,
	) (i InstanceID, ok bool, err error)

	// SaveKeys persists the changes to a saga instance's mapping key set.
	//
	// sn is the name of the saga.
	SaveKeys(
		ctx context.Context,
		tx persistence.Tx,
		sn string,
		id InstanceID,
		ks KeySet,
	) error
}

// KeySet is a set of "mapping keys" that are associated with a saga instance.
type KeySet map[string]struct{}

// NewKeySet returns a key set containing the given keys.
func NewKeySet(keys ...string) KeySet {
	s := make(KeySet, len(keys))

	for _, k := range keys {
		if k == "" {
			panic("mapping keys must not be empty")
		}

		s[k] = struct{}{}
	}

	return s
}
