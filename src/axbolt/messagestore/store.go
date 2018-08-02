package messagestore

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/messagestore"
	"github.com/jmalloc/ax/src/ax/persistence"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

// StreamBktName is the name of the Bolt root bucket that contains message
// offsets for each stream
var StreamBktName = []byte("ax_messagestore_stream")

// MessageBktName is the name of the Bolt root bucket that contains all messages
// stored as a global stream
var MessageBktName = []byte("ax_messagestore_message")

// offset key is the key inside stream buckets to retain the value of tha latesgt
// offset used inside the stream bucket.
var offsetkey = []byte("offset")

// msgbkt is the key inside the stream bucket to address a message subbucket.
var msgbkt = []byte("msgs")

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
	return &Stream{
		Fetcher: &StreamFetcher{
			DB:     boltpersistence.ExtractDB(ds),
			Stream: stream,
		},
		NextOffset: offset,
	}, true, nil
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
