package direct

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/saga"
)

// Saga is an interface for sagas that use key set mapping.
type Saga interface {
	saga.Saga

	// InstanceIDForMessage returns the ID of the saga instance to which the
	// given message is routed, if any.
	//
	// If ok is false the message is ignored; otherwise, the message is routed
	// to the saga instance with the returned ID.
	InstanceIDForMessage(ctx context.Context, env ax.Envelope) (saga.InstanceID, bool, error)
}
