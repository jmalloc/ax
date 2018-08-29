package messagestore

import (
	"context"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/messagestore"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axbolt/internal/boltutil"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

const (
	// streamBktName is the name of the Bolt root bucket that contains message
	// offsets for each stream
	streamBktName = "ax_messagestore_stream"

	// globalStreamBktName is the name of the Bolt root bucket that contains all messages
	// stored as a global stream
	globalStreamBktName = "ax_messagestore_message"

	// offsetKey is the key within stream buckets to hold the value of the latest
	// stream offset. This value is incremented in case of successful message
	// insertion.
	offsetKey = "offset"

	// msgsBktName is the name of the bucket to store the messages.
	msgsBktName = "msgs"
)

// Store is a Bolt-backed implementation of Ax's
// messagestore.GloballyOrderedStore interface.
type Store struct{}

// AppendMessages appends one or more messages to a named stream.
//
// offset is a zero-based index into the stream. An error is returned if
// offset is not the next unused offset in the stream.
func (Store) AppendMessages(
	ctx context.Context,
	ptx persistence.Tx,
	stream string,
	offset uint64,
	envs []ax.Envelope,
) error {
	tx := boltpersistence.ExtractTx(ptx)
	for _, env := range envs {
		global, err := insertGlobalOffsetMessage(tx, env)
		if err != nil {
			return err
		}
		if err = insertStreamOffset(tx, stream, offset, global); err != nil {
			return err
		}
		offset++
	}
	return nil
}

// OpenStream opens a stream of messages for reading from a specific offset.
//
// The offset may be past the end of the stream. It returns false if the stream
// does not exist.
func (Store) OpenStream(
	ctx context.Context,
	ds persistence.DataStore,
	stream string,
	offset uint64,
) (messagestore.Stream, bool, error) {
	db := boltpersistence.ExtractDB(ds)
	ok, err := streamExists(db, stream)
	if err != nil {
		return nil, ok, err
	}

	return &Stream{
		Fetcher: &StreamFetcher{
			DB:     db,
			Stream: stream,
		},
		NextOffset: offset,
	}, ok, nil
}

// OpenGlobal opens the entire store for reading as a single stream.
//
// The offset may be beyond the end of the stream.
func (Store) OpenGlobal(
	ctx context.Context,
	ds persistence.DataStore,
	offset uint64,
) (messagestore.Stream, error) {
	return &Stream{
		Fetcher: &GlobalFetcher{
			DB: boltpersistence.ExtractDB(ds),
		},
		NextOffset: offset,
	}, nil
}

// streamExists returns true if the stream exists in the store.
// It returns an error if starting db transaction fails
func streamExists(db *bolt.DB, s string) (bool, error) {
	tx, err := db.Begin(false)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	bkt := boltutil.GetBkt(tx, streamBktName, s)
	return bkt != nil, nil
}
