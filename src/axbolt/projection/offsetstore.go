package projection

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axbolt/internal/boltutil"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

const (
	// projectionOffsetBktName is the name of the Bolt root bucket where
	// projection offset data is stored.
	projectionOffsetBktName = "ax_projection_offset"
)

// OffsetStore is a Bolt-backed implementation of Ax's projection.OffsetStore
// interface.
type OffsetStore struct{}

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

	if b := boltutil.Get(tx, pk, projectionOffsetBktName); b != nil {
		return binary.BigEndian.Uint64(b), nil
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
	tx := boltpersistence.ExtractTx(ptx)

	if c != 0 {
		b := boltutil.Get(tx, pk, projectionOffsetBktName)
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

	b := [8]byte{}
	c++
	binary.BigEndian.PutUint64(b[:], c)
	return boltutil.Put(tx, pk, b[:], projectionOffsetBktName)
	// TODO: use OCC error https://github.com/jmalloc/ax/issues/93
}
