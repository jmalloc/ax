package keyset

import (
	"context"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/saga"
)

// Resolver is an interface that provides the application-defined logic for
// mapping a message to its target saga instance.
type Resolver interface {
	// GenerateInstanceID returns the saga ID to use for a new instance.
	//
	// It is called when a "trigger" message is received and there is no
	// existing saga instance. env contains the "trigger" message.
	GenerateInstanceID(ctx context.Context, env ax.Envelope) (id saga.InstanceID, err error)

	// MappingKeyForMessage returns the key used to locate the saga instance
	// to which the given message is routed, if any.
	//
	// If ok is false the message is ignored; otherwise, the message is routed
	// to the saga instance that contains k in its associated key set.
	//
	// New saga instances are created when no matching instance can be found
	// and the message is declared as a "trigger" by the saga's MessageTypes()
	// method; otherwise, HandleNotFound() is called.
	MappingKeyForMessage(ctx context.Context, env ax.Envelope) (k string, ok bool, err error)

	// MappingKeysForInstance returns the set of mapping keys associated with
	// the given instance.
	//
	// When a message is received, a mapping key is produced by calling
	// MappingKeyForMessage(). The message is routed to the saga instance that
	// contains this key in its key set.
	//
	// Key sets must be disjoint. That is, no two instances of the same saga
	// may share any keys.
	//
	// The key set is rebuilt whenever a message is received. It is persisted
	// alongside the saga instance by the Repository.
	//
	// Extra care must be taken when introducing a new key to the set, as the key
	// set for existing saga instances will not be updated until they next receive
	// a message.
	MappingKeysForInstance(context.Context, saga.Instance) ([]string, error)
}
