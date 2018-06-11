package direct

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// Mapper is an implementation of saga.Mapper that maps messages to saga
// instances by having the saga implement a method that returns the instance ID
// directly.
//
// The saga must implement the direct.Saga interface to use direct mapping.
type Mapper struct{}

// MapMessageToInstance returns the ID of the saga instance that is the target
// of the given message.
//
// It returns false if the message should be ignored.
func (m *Mapper) MapMessageToInstance(
	ctx context.Context,
	sg saga.Saga,
	tx persistence.Tx,
	env ax.Envelope,
) (saga.InstanceID, bool, error) {
	return sg.(Saga).InstanceIDForMessage(ctx, env)
}

// UpdateMapping notifies the mapper that an instance has been modified,
// allowing it to update it's mapping information, if necessary.
func (m *Mapper) UpdateMapping(
	ctx context.Context,
	sg saga.Saga,
	tx persistence.Tx,
	i saga.Instance,
) error {
	return nil
}

// DeleteMapping notifies the mapper that an instance has been completed,
// allowing it to remove it's mapping information, if necessary.
func (m *Mapper) DeleteMapping(
	ctx context.Context,
	sg saga.Saga,
	tx persistence.Tx,
	i saga.Instance,
) error {
	return nil
}
