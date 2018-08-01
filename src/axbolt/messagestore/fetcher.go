package messagestore

import (
	"bytes"
	"context"
	"encoding/binary"

	bolt "github.com/coreos/bbolt"
	"github.com/golang/protobuf/proto"
)

// Fetcher is an interface for fetching rows from the message store
type Fetcher interface {
	// FetchRows fetches the n rows beginning at the given offset.
	FetchRows(
		ctx context.Context,
		offset,
		n uint64,
	) (map[uint64]*StoredMessage, error)
}

// StreamFetcher is a fetcher that fetches rows for a specific stream.
type StreamFetcher struct {
	DB     *bolt.DB
	Stream string
}

// FetchRows fetches the n rows beginning at the given offset.
func (f *StreamFetcher) FetchRows(
	ctx context.Context,
	offset,
	n uint64,
) (map[uint64]*StoredMessage, error) {
	tx, err := f.DB.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bkt, err := tx.CreateBucketIfNotExists(MessageBktName)
	if err != nil {
		return nil, err
	}
	if bkt, err = bkt.CreateBucketIfNotExists([]byte(f.Stream)); err != nil {
		return nil, err
	}
	if bkt, err = bkt.CreateBucketIfNotExists(msgbkt); err != nil {
		return nil, err
	}
	gbkt, err := tx.CreateBucketIfNotExists(MessageBktName)
	if err != nil {
		return nil, err
	}

	c := bkt.Cursor()
	r := map[uint64]*StoredMessage{}
	b := [8]byte{}

	binary.BigEndian.PutUint64(b[:], offset)
	k, v := c.Seek(b[:])
	// no exact match, return an empty map
	if k == nil || !bytes.Equal(k, b[:]) {
		return r, nil
	}
	// iterate over stream bucket and retrieve n number of messages from global
	// bucket using v as a global offset, after that insert them into resultant map
	// with k as a stream offset value. If a message for whatever reason doesn't
	// exist in global offset bucket, it is skipped. The number of messages in the
	// bucket can be less than n.
	for i := uint64(0); i < n && k != nil && v != nil; k, v = c.Next() {
		if m := gbkt.Get(v); m != nil {
			var msg StoredMessage
			if err = proto.Unmarshal(m, &msg); err != nil {
				return nil, err
			}
			r[binary.BigEndian.Uint64(k)] = &msg
		}
		i++
	}
	return r, nil
}

// GlobalFetcher is a fetcher that fetches rows for the entire store
type GlobalFetcher struct {
	DB *bolt.DB
}

// FetchRows fetches the n rows beginning at the given offset.
func (f *GlobalFetcher) FetchRows(
	ctx context.Context,
	offset, n uint64,
) (map[uint64]*StoredMessage, error) {
	tx, err := f.DB.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bkt, err := tx.CreateBucketIfNotExists(MessageBktName)
	if err != nil {
		return nil, err
	}

	c := bkt.Cursor()
	r := map[uint64]*StoredMessage{}
	b := [8]byte{}

	binary.BigEndian.PutUint64(b[:], offset)
	k, v := c.Seek(b[:])
	// no exact match, return an empty map
	if k == nil || !bytes.Equal(k, b[:]) {
		return r, nil
	}
	// iterate over the global offset bucket and retrieve n number of messages
	// starting from the message with key equal to the value of offset. The number
	// of messages in the bucket can be less than n.
	for i := uint64(0); i < n && k != nil && v != nil; k, v = c.Next() {
		var msg StoredMessage
		if err = proto.Unmarshal(v, &msg); err != nil {
			return nil, err
		}
		r[binary.BigEndian.Uint64(k)] = &msg
		i++
	}
	return r, nil
}
