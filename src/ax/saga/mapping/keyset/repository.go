package keyset

import (
	"context"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// Repository is an interface for storing and querying saga instances
// "key sets".
type Repository interface {
	// FindByKey returns the ID of a saga instance that has a specific key in
	// its key set.
	//
	// pk is the saga's persistence key, mk is the mapping key.
	// ok is false if no saga instance has a key set containing mk.
	FindByKey(
		ctx context.Context,
		tx persistence.Tx,
		pk, mk string,
	) (i saga.InstanceID, ok bool, err error)

	// SaveKeys associates a set of mapping keys with a saga instance.
	//
	// Key sets must be disjoint. That is, no two instances of the same saga
	// may share any keys.
	//
	// pk is the saga's persistence key. ks is the set of mapping keys.
	//
	// SaveKeys() may panic if ks contains duplicate keys.
	SaveKeys(
		ctx context.Context,
		tx persistence.Tx,
		pk string,
		ks []string,
		id saga.InstanceID,
	) error

	// DeleteKeys removes any mapping keys associated with a saga instance.
	//
	// pk is the saga's persistence key.
	DeleteKeys(
		ctx context.Context,
		tx persistence.Tx,
		pk string,
		id saga.InstanceID,
	) error
}
