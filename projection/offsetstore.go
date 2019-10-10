package projection

import (
	"context"

	"github.com/jmalloc/ax/persistence"
)

// OffsetStore is an interface for persisting a consumer's current position in a
// message stream.
type OffsetStore interface {
	// LoadOffset returns the offset at which a consumer should resume
	// reading from the stream.
	//
	// pk is the projector's persistence key.
	LoadOffset(
		ctx context.Context,
		ds persistence.DataStore,
		pk string,
	) (uint64, error)

	// IncrementOffset increments the offset at which a consumer should resume
	// reading from the stream by one.
	//
	// pk is the projector's persitence key. c is the offset that is currently
	// stored, as returned by LoadOffset(). If c is not the offset that is
	// currently stored, the increment fails and a non-nil error is returned.
	IncrementOffset(
		ctx context.Context,
		tx persistence.Tx,
		pk string,
		c uint64,
	) error
}
