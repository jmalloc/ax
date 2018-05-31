package saga

import (
	"context"
	"errors"
	"fmt"

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
	//
	// SaveKeys() may panic if ks contains duplicate keys.
	SaveKeys(
		ctx context.Context,
		tx persistence.Tx,
		sn string,
		id InstanceID,
		ks []string,
	) error
}

// validateKeySet returns a copy of ks with duplicate keys removed.
// It returns an error of any of the keys is the empty string.
func validateKeySet(ks []string) ([]string, error) {
	dedup := make([]string, 0, len(ks))
	seen := make(map[string]struct{}, len(ks))

	for _, k := range ks {
		if k == "" {
			return nil, errors.New("mapping keys must not be empty")
		}

		if _, ok := seen[k]; ok {
			return nil, fmt.Errorf("the mapping key %s is repeated in the key set", k)
		}

		seen[k] = struct{}{}
		dedup = append(dedup, k)
	}

	return dedup, nil
}
