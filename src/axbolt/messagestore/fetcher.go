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
	// FetchMessages fetches the n messages beginning at the given offset.
	FetchMessages(
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

// FetchMessages fetches the n messages beginning at the given offset.
// If no messages available for the stream, this method returns an empty map.
// If any errors occur in the process of fetching messages, nil and error object
// are returned.
func (f *StreamFetcher) FetchMessages(
	ctx context.Context,
	offset,
	n uint64,
) (map[uint64]*StoredMessage, error) {
	type res struct {
		msgs map[uint64]*StoredMessage
		err  error
	}
	resNotify := make(chan res)
	tx, err := f.DB.Begin(false)
	if err != nil {
		return nil, err
	}

	go func() {
		msgs, err := f.fetch(ctx, tx, offset, n)
		resNotify <- res{
			msgs: msgs,
			err:  err,
		}
	}()

	select {
	case <-ctx.Done():
		_ = tx.Rollback()
		return nil, ctx.Err()
	case r := <-resNotify:
		return r.msgs, r.err
	}
}

// fetch iterates over a stream bucket and retrieve n number of messages from global
// bucket using v as a global offset, after that insert them into resultant map
// with k as a stream offset value. If a message for whatever reason doesn't
// exist in global offset bucket, it is skipped. The number of messages in the
// bucket can be less than n.
func (f StreamFetcher) fetch(
	ctx context.Context,
	tx *bolt.Tx,
	offset,
	n uint64,
) (map[uint64]*StoredMessage, error) {
	// check if context is already cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	r := map[uint64]*StoredMessage{}
	bkt := tx.Bucket(StreamBktName)
	if bkt == nil {
		return r, nil
	}
	if bkt = bkt.Bucket([]byte(f.Stream)); bkt == nil {
		return r, nil
	}
	if bkt = bkt.Bucket(msgbkt); bkt == nil {
		return r, nil
	}
	gbkt := tx.Bucket(MessageBktName)
	if gbkt == nil {
		return r, nil
	}
	if gbkt = gbkt.Bucket(msgbkt); gbkt == nil {
		return r, nil
	}

	c := bkt.Cursor()
	b := [8]byte{}

	binary.BigEndian.PutUint64(b[:], offset)
	k, v := c.Seek(b[:])
	// no exact match or value greater than offset, return an empty map
	if k == nil || bytes.Compare(k, b[:]) == -1 {
		return r, nil
	}

	for i := uint64(0); i < n && k != nil && v != nil; k, v = c.Next() {
		if m := gbkt.Get(v); m != nil {
			var msg StoredMessage
			if err := proto.Unmarshal(m, &msg); err != nil {
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

// FetchMessages fetches the n messages beginning at the given offset.
func (f *GlobalFetcher) FetchMessages(
	ctx context.Context,
	offset, n uint64,
) (map[uint64]*StoredMessage, error) {
	tx, err := f.DB.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	r := map[uint64]*StoredMessage{}

	bkt := tx.Bucket(MessageBktName)
	if bkt == nil {
		return r, nil
	}
	if bkt = bkt.Bucket(msgbkt); bkt == nil {
		return r, nil
	}

	c := bkt.Cursor()
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
