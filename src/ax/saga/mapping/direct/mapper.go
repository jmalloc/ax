package direct

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// Mapper is an implementation of saga.Mapper that maps messages to saga
// by treating a particular field in the message as the instance ID.
//
// The saga must implement the direct.Saga interface to use direct mapping.
type Mapper struct{}

// MapMessageToInstance returns the ID of the saga instance that is the target
// of the given message.
//
// If no existing saga instance is found, it returns false.
func (m *Mapper) MapMessageToInstance(
	ctx context.Context,
	sg saga.Saga,
	tx persistence.Tx,
	env ax.Envelope,
) (saga.MapResult, saga.InstanceID, error) {
	id, ok, err := sg.(Saga).InstanceIDForMessage(ctx, env)
	if !ok || err != nil {
		return saga.MapResultIgnore, saga.InstanceID{}, err
	}

	return saga.MapResultFound, id, nil
}

// UpdateMapping notifies the mapper that a message has been handled by
// an instance. Giving it the oppurtunity to update mapping data to reflect
// the changes, if necessary.
func (m *Mapper) UpdateMapping(
	ctx context.Context,
	sg saga.Saga,
	tx persistence.Tx,
	i saga.Instance,
) error {
	return nil
}
