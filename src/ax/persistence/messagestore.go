package persistence

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// MessageStore is an interface for persisting streams of messages.
type MessageStore interface {
	// AppendMessages appends one or more messages to a named stream.
	//
	// offset is a zero-based index into the stream. An error is returned if
	// offset is not the next unused offset in the stream.
	AppendMessages(
		ctx context.Context,
		tx Tx,
		stream string,
		offset uint64,
		envs []ax.Envelope,
	) error

	// OpenStream opens a stream of messages for reading from a specific offset.
	//
	// The offset may be past the end of the stream. It returns an error if
	// the stream does not exist.
	OpenStream(
		ctx context.Context,
		tx Tx,
		stream string,
		offset uint64,
	) (MessageStream, error)
}

// MessageStream is a stream of messages stored in a MessageStore.
type MessageStream interface {
	// Next advances the stream to the next message.
	// It returns false if there are no more messages in the stream.
	Next(ctx context.Context) (bool, error)

	// Get returns the message at the current offset in the stream.
	Get(ctx context.Context) (ax.Envelope, error)

	// Close closes the stream.
	Close() error
}
