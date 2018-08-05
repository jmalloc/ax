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

	if o := boltutil.Get(tx, OffsetKey, GlobalStreamBktName); o != nil {
		offset = binary.BigEndian.Uint64(o)
	}

	bkt, err := boltutil.MakeBktWithPath(tx, GlobalStreamMsgBktPath)
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
	if err = boltutil.MarshalProto(bkt, b[:], m); err != nil {
		return 0, err
	}

	return offset,
		boltutil.Put(tx, OffsetKey, b[:], GlobalStreamBktName)
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
	p := fmt.Sprintf("%s/%s", StreamBktName, stream)
	o := boltutil.GetWithPath(tx, OffsetKey, p)
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

	bkt, err := boltutil.MakeBktWithPath(tx, fmt.Sprintf("%s/msgs", p))
	if err != nil {
		return err
	}
	if err := bkt.Put(b1[:], b2[:]); err != nil {
		return err
	}
	return boltutil.PutWithPath(tx, OffsetKey, b1[:], p)
}
