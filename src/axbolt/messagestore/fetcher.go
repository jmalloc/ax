package messagestore

import (
	"bytes"
	"context"
	"encoding/binary"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/ax/src/axbolt/internal/boltutil"
)

// Fetcher is an interface for fetching rows from the message store
type Fetcher interface {
	// FetchMessages fetches the n messages beginning at the given offset.
	FetchMessages(
		ctx context.Context,
		offset,
		n uint64,
	) (map[uint64]*ax.EnvelopeProto, error)
}

// StreamFetcher is a fetcher that fetches rows for a specific stream.
type StreamFetcher struct {
	DB     *bolt.DB
	Stream string
}

// FetchMessages fetches the n messages beginning at the given offset. If no
// messages available for the stream, this method returns an empty map. If any
// errors occur in the process of fetching messages, nil and error object are
// returned.
func (f *StreamFetcher) FetchMessages(
	ctx context.Context,
	offset,
	n uint64,
) (map[uint64]*ax.EnvelopeProto, error) {
	type res struct {
		envpbs map[uint64]*ax.EnvelopeProto
		err    error
	}
	resNotify := make(chan res)
	tx, err := f.DB.Begin(false)
	if err != nil {
		return nil, err
	}

	go func() {
		envpbs, err := f.fetch(ctx, tx, offset, n)
		resNotify <- res{
			envpbs: envpbs,
			err:    err,
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-resNotify:
		return r.envpbs, r.err
	}
}

// fetch iterates over a stream bucket and retrieves n number of messages from
// global bucket using v as a global offset, after that insert them into
// resultant map with k as a stream offset value. If a message for whatever
// reason doesn't exist in global offset bucket, it is skipped. The number of
// messages in the bucket can be less than n.
func (f StreamFetcher) fetch(
	ctx context.Context,
	tx *bolt.Tx,
	offset,
	n uint64,
) (map[uint64]*ax.EnvelopeProto, error) {
	// check if context is already canceled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r := map[uint64]*ax.EnvelopeProto{}
	bkt := boltutil.GetBkt(
		tx,
		streamBktName,
		f.Stream,
		msgsBktName,
	)
	if bkt == nil {
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
		if m := boltutil.Get(
			tx,
			string(v),
			globalStreamBktName,
			msgsBktName,
		); m != nil {

			var envpb ax.EnvelopeProto
			if err := proto.Unmarshal(m, &envpb); err != nil {
				return r, nil
			}
			r[binary.BigEndian.Uint64(k)] = &envpb
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
) (map[uint64]*ax.EnvelopeProto, error) {
	type res struct {
		envpbs map[uint64]*ax.EnvelopeProto
		err    error
	}
	resNotify := make(chan res)
	tx, err := f.DB.Begin(false)
	if err != nil {
		return nil, err
	}

	go func() {
		envpbs, err := f.fetch(ctx, tx, offset, n)
		resNotify <- res{
			envpbs: envpbs,
			err:    err,
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-resNotify:
		return r.envpbs, r.err
	}
}

// fetch iterates over the global offset bucket and retrieves n number of
// messages starting from the message with key equal to the value of offset. The
// number of messages in the bucket can be less than n.
func (f *GlobalFetcher) fetch(
	ctx context.Context,
	tx *bolt.Tx,
	offset,
	n uint64,
) (map[uint64]*ax.EnvelopeProto, error) {
	// check if context is already canceled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	r := map[uint64]*ax.EnvelopeProto{}
	bkt := boltutil.GetBkt(tx, globalStreamBktName, msgsBktName)
	if bkt == nil {
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
		var envpb ax.EnvelopeProto
		if err := proto.Unmarshal(v, &envpb); err != nil {
			return r, err
		}
		r[binary.BigEndian.Uint64(k)] = &envpb
		i++
	}
	return r, nil
}
