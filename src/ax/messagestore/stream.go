package messagestore

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Stream is an interface for reading an ordered stream of messages.
type Stream interface {
	// TryNext advances the stream to the next message.
	//
	// It returns false if there are no more messages in the stream.
	TryNext(ctx context.Context) (bool, error)

	// Get returns the message at the current offset in the stream.
	Get(ctx context.Context) (ax.Envelope, error)

	// Offset returns the offset of the message returned by Get().
	Offset() (uint64, error)

	// Close closes the stream.
	Close() error
}
