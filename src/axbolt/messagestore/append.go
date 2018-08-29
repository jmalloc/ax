package messagestore

import (
	"encoding/binary"
	fmt "fmt"

	"github.com/jmalloc/ax/src/axbolt/internal/boltutil"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/ax/src/ax"
)

// insertGlobalOffsetMessage inserts a message into the global offset bucket.
// This method retrieves the previous global offset and inserts the new message
// with an incremented offset.
//
// If the bucket does not have any previous keys it will insert the message with
// offset 1. It returns an error if it occurs in the process of creating offset
// bucket (in case does not exist) marchalling and putting the message into bkt
// bucket.
func insertGlobalOffsetMessage(tx *bolt.Tx, env ax.Envelope) (uint64, error) {
	var (
		offset uint64
		b      [8]byte
	)

	if o := boltutil.Get(tx, offsetKey, globalStreamBktName); o != nil {
		offset = binary.BigEndian.Uint64(o)
	}

	envb, err := ax.MarshalEnvelope(env)
	if err != nil {
		return 0, err
	}

	offset++
	binary.BigEndian.PutUint64(b[:], offset)
	if err = boltutil.PutB(
		tx,
		b[:],
		envb,
		globalStreamBktName,
		msgsBktName,
	); err != nil {
		return 0, err
	}

	return offset,
		boltutil.Put(tx, offsetKey, b[:], globalStreamBktName)
}

// insertStreamOffset inserts a message into the stream offset bucket. This
// method retrieves the previous stream offset and inserts the new message with
// an incremented offset.
//
// If the bucket does not have any previous keys it will insert the message with
// offset 1. It returns error if offset does not equal the latest offset of tne steam.
// It returns an error if it occurs in the process of creating offset
// bucket (in case does not exist) marchalling and putting the message into the
// bolt bucket.
func insertStreamOffset(
	tx *bolt.Tx,
	stream string,
	offset, global uint64,
) error {
	var (
		b1, b2 [8]byte
	)
	o := boltutil.Get(tx, offsetKey, streamBktName, stream)
	if o != nil && binary.BigEndian.Uint64(o) != offset ||
		o == nil && offset != 0 {
		// TODO: use OCC error https://github.com/jmalloc/ax/issues/93
		return fmt.Errorf(
			"can not append to stream %s, %d is not the next free offset",
			stream,
			offset,
		)
	}

	offset++
	binary.BigEndian.PutUint64(b1[:], offset)
	binary.BigEndian.PutUint64(b2[:], global)

	if err := boltutil.PutB(
		tx,
		b1[:],
		b2[:],
		streamBktName, stream, msgsBktName,
	); err != nil {
		return err
	}

	return boltutil.Put(tx, offsetKey, b1[:], streamBktName, stream)
}
