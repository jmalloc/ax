package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

// KeySetRepository is a MySQL-backed implementation of Ax's keyset.Repository
// interface.
type KeySetRepository struct{}

// FindByKey returns the ID of a saga instance that has a specific key in
// its key set.
//
// pk is the saga's persistence key, mk is the mapping key.
// ok is false if no saga instance has a key set containing mk.
func (KeySetRepository) FindByKey(
	ctx context.Context,
	ptx persistence.Tx,
	pk, mk string,
) (id saga.InstanceID, ok bool, err error) {
	tx := boltpersistence.ExtractTx(ptx)

	bkt := tx.Bucket([]byte("ax_saga"))
	if bkt == nil {
		return
	}

	bkt = bkt.Bucket([]byte("ax_saga_keyset"))
	if bkt == nil {
		return
	}

	bkt = bkt.Bucket([]byte(pk))
	if bkt == nil {
		return
	}

	raw := bkt.Get([]byte(mk))
	if raw == nil {
		return
	}

	err = id.Parse(string(raw))
	if err != nil {
		return
	}
	ok = true
	return
}

// SaveKeys associates a set of mapping keys with a saga instance.
//
// Key sets must be disjoint. That is, no two instances of the same saga
// may share any keys.
//
// pk is the saga's persistence key. ks is the set of mapping keys.
//
// SaveKeys() may panic if ks contains duplicate keys.
func (KeySetRepository) SaveKeys(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	ks []string,
	id saga.InstanceID,
) error {
	tx := boltpersistence.ExtractTx(ptx)

	bkt, err := tx.CreateBucketIfNotExists([]byte("ax_saga"))
	if err != nil {
		return err
	}
	bkt, err = bkt.CreateBucketIfNotExists([]byte("ax_saga_keyset"))
	if err != nil {
		return err
	}
	bkt, err = bkt.CreateBucketIfNotExists([]byte(pk))
	if err != nil {
		return err
	}

	c := bkt.Cursor()
	if k, v := c.First(); k != nil && v != nil {
		for k, _ = c.Last(); k != nil; k, _ = c.Prev() {
			if err = c.Delete(); err != nil {
				return err
			}
		}
	}

	for _, k := range ks {
		if err = bkt.Put([]byte(k), []byte(id.Get())); err != nil {
			return err
		}
	}

	return nil
}
