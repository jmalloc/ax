package messagestore

import (
	"context"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/persistence"
)

// Store is an interface for manipulating persisted streams of messages.
type Store interface {
	// AppendMessages appends one or more messages to a named stream.
	//
	// offset is a zero-based index into the stream. An error is returned if
	// offset is not the next unused offset in the stream.
	AppendMessages(
		ctx context.Context,
		tx persistence.Tx,
		stream string,
		offset uint64,
		envs []ax.Envelope,
	) error

	// OpenStream opens a stream of messages for reading from a specific offset.
	//
	// The offset may be beyond the end of the stream. It returns false if the
	// stream does not exist.
	OpenStream(
		ctx context.Context,
		ds persistence.DataStore,
		stream string,
		offset uint64,
	) (Stream, bool, error)
}

// GloballyOrderedStore is a store that preserves a global ordering for messages
// across all streams, allowing the entire store to be consumed via a single
// stream.
type GloballyOrderedStore interface {
	Store

	// OpenGlobal opens the entire store for reading as a single stream.
	//
	// The offset may be beyond the end of the stream.
	OpenGlobal(
		ctx context.Context,
		ds persistence.DataStore,
		offset uint64,
	) (Stream, error)
}
