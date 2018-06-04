package messagestore

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
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
	// The offset may be past the end of the stream. It returns false if the
	// stream does not exist.
	OpenStream(
		ctx context.Context,
		ds persistence.DataStore,
		stream string,
		offset uint64,
	) (Stream, bool, error)
}
