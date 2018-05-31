package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax/persistence"
)

// KeySetRepository is an interface for storing and querying saga instances
// "key sets".
type KeySetRepository interface {
	// FindByKey returns the ID of the saga instance that contains k in its
	// key set for the saga named sn.
	//
	// ok is false if no saga instance has a key set containing k.
	FindByKey(
		ctx context.Context,
		tx persistence.Tx,
		sn, k string,
	) (i InstanceID, ok bool, err error)

	// SaveKeys associates a key set with the saga instance identified by id
	// for the saga named sn.
	//
	// Key sets must be disjoint. That is, no two instances of the same saga
	// may share any keys.
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
