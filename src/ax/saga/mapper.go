package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// Mapper is an interface for mapping inbound messages to their target saga
// instance.
type Mapper interface {
	// MapMessageToInstance returns the ID of the saga instance that is the
	// target of the given message.
	//
	// It returns false if the message should be ignored.
	MapMessageToInstance(
		ctx context.Context,
		sg Saga,
		tx persistence.Tx,
		env ax.Envelope,
	) (InstanceID, bool, error)

	// UpdateMapping notifies the mapper that an instance has been modified,
	// allowing it to update it's mapping information, if necessary.
	UpdateMapping(
		ctx context.Context,
		sg Saga,
		tx persistence.Tx,
		i Instance,
	) error

	// DeleteMapping notifies the mapper that an instance has been completed,
	// allowing it to remove it's mapping information, if necessary.
	DeleteMapping(
		ctx context.Context,
		sg Saga,
		tx persistence.Tx,
		i Instance,
	) error
}
