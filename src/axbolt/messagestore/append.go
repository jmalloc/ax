package messagestore

import (
	"encoding/binary"
	fmt "fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
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
	bl0, err := tx.CreateBucketIfNotExists(MessageBktName)
	if err != nil {
		return 0, err
	}

	if o := bl0.Get(offsetkey); o != nil {
		offset = binary.BigEndian.Uint64(o)
	}

	bl1, err := bl0.CreateBucketIfNotExists(msgbkt)
	if err != nil {
		return 0, err
	}

	m := &StoredMessage{
		Id:            env.MessageID.Get(),
		CausationId:   env.CausationID.Get(),
		CorrelationId: env.CorrelationID.Get(),
		CreatedAt:     env.CreatedAt.Format(time.RFC3339Nano),
		SendAt:        env.SendAt.Format(time.RFC3339Nano),
	}
	m.Message, err = ptypes.MarshalAny(env.Message)
	if err != nil {
		return 0, err
	}

	offset++
	binary.BigEndian.PutUint64(b[:], offset)
	if err = boltutil.MarshalProto(bl1, b[:], m); err != nil {
		return 0, err
	}
	return offset, bl0.Put(offsetkey, b[:])
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
	bl0, err := tx.CreateBucketIfNotExists(StreamBktName)
	if err != nil {
		return err
	}

	bl1, err := bl0.CreateBucketIfNotExists([]byte(stream))
	if err != nil {
		return err
	}

	if o := bl1.Get(offsetkey); offset != binary.BigEndian.Uint64(o) ||
		(o == nil && offset != 0) {
		// TODO: use OCC error https://github.com/jmalloc/ax/issues/93
		return fmt.Errorf(
			"can not append to stream %s, %d is not the next free offset",
			stream,
			offset,
		)
	}

	bl2, err := bl1.CreateBucketIfNotExists(msgbkt)
	if err != nil {
		return err
	}

	offset++
	binary.BigEndian.PutUint64(b1[:], offset)
	binary.BigEndian.PutUint64(b2[:], global)
	if err := bl2.Put(b1[:], b2[:]); err != nil {
		return err
	}
	return bl1.Put(offsetkey, b1[:])
}
