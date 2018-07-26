package projection

import (
	"context"
	"encoding/binary"
	"fmt"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/ax/src/ax/persistence"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

// OffsetStore is a Bolt-backed implementation of Ax's projection.OffsetStore
// interface.
type OffsetStore struct{}

// ProjectionOffsetBktName is the name of of the Bolt root bucket where
// projection offset data is stored.
var ProjectionOffsetBktName = []byte("ax_projection_offset")

// LoadOffset returns the offset at which a consumer should resume
// reading from the stream.
//
// pk is the projector's persistence key.
func (OffsetStore) LoadOffset(
	ctx context.Context,
	ds persistence.DataStore,
	pk string,
) (uint64, error) {
	db := boltpersistence.ExtractDB(ds)
	tx, err := db.Begin(false)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if bkt := tx.Bucket(ProjectionOffsetBktName); bkt != nil {
		if b := bkt.Get([]byte(pk)); b != nil {
			return binary.BigEndian.Uint64(b), nil
		}
	}
	return 0, nil
}

// IncrementOffset increments the offset at which a consumer should resume
// reading from the stream by one.
//
// pk is the projector's persistence key. c is the offset that is currently
// stored, as returned by LoadOffset(). If c is not the offset that is
// currently stored, the increment fails and a non-nil error is returned.
func (s OffsetStore) IncrementOffset(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	c uint64,
) error {

	var (
		bkt *bolt.Bucket
		err error
	)
	tx := boltpersistence.ExtractTx(ptx)
	if bkt, err = tx.CreateBucketIfNotExists(
		ProjectionOffsetBktName,
	); err != nil {
		return err
	}

	if c != 0 {
		b := bkt.Get([]byte(pk))
		if b == nil {
			return nil
		}
		offset := binary.BigEndian.Uint64(b)
		if c != offset {
			return fmt.Errorf(
				"can not increment projection offset for persistence key %s, offset %d is not the current offset",
				pk,
				c,
			)
		}
	}

	c++
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, c)
	return bkt.Put([]byte(pk), b)
	// TODO: use OCC error https://github.com/jmalloc/ax/issues/93
}
