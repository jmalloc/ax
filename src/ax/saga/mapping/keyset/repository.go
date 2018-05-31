package keyset

import (
	"context"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// Repository is an interface for storing and querying saga instances
// "key sets".
type Repository interface {
	// FindByKey returns the ID of the saga instance that contains k in its
	// key set for the saga named sn.
	//
	// ok is false if no saga instance has a key set containing k.
	FindByKey(
		ctx context.Context,
		tx persistence.Tx,
		sn, k string,
	) (i saga.InstanceID, ok bool, err error)

	// SaveKeys associates a key set with the saga instance identified by id
	// for the saga named sn.
	//
	// Key sets must be disjoint. That is, no two instances of the same saga
	// may share any keys.
	//
	// SaveKeys() may panic if ks contains duplicate keys.
	SaveKeys(
		ctx context.Context,
		tx persistence.Tx,
		sn string,
		id saga.InstanceID,
		ks []string,
	) error
}
