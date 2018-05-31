package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// MapResult contains the result of a call to Mapper.MapMessageToInstance()
type MapResult int

const (
	// MapResultFound indicates that an instance ID was found for the supplied
	// message.
	MapResultFound MapResult = iota

	// MapResultNotFound indicates that no existing instance ID was found for the
	// supplied message, and that a new instance should be created if possible.
	MapResultNotFound

	// MapResultIgnore indicates that this message should not be routed to any
	// instance.
	MapResultIgnore
)

// Mapper is an interface for mapping inbound messages to their target saga
// instance.
type Mapper interface {
	// MapMessageToInstance returns the ID of the saga instance that is the
	// target of the given message.
	//
	// If no existing saga instance is found, it returns false.
	MapMessageToInstance(
		ctx context.Context,
		sg Saga,
		tx persistence.Tx,
		env ax.Envelope,
	) (MapResult, InstanceID, error)

	// UpdateMapping notifies the mapper that a message has been handled by
	// an instance. Giving it the oppurtunity to update mapping data to reflect
	// the changes, if necessary.
	UpdateMapping(
		ctx context.Context,
		sg Saga,
		tx persistence.Tx,
		i Instance,
	) error
}
