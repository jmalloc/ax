package direct

import (
	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/saga"
)

// Resolver is an interface that provides the application-defined logic for
// mapping a message to its target saga instance.
type Resolver interface {
	// InstanceIDForMessage returns the ID of the saga instance to which the
	// given message is routed, if any.
	//
	// If ok is false the message is ignored; otherwise, the message is routed
	// to the saga instance with the returned ID.
	InstanceIDForMessage(env ax.Envelope) (saga.InstanceID, bool)
}
