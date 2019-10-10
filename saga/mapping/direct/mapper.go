package direct

import (
	"context"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/persistence"
	"github.com/jmalloc/ax/saga"
)

// Mapper is an implementation of saga.Mapper that maps messages to saga
// instances using only information contained in the message.
type Mapper struct {
	Resolver Resolver
}

// MapMessageToInstance returns the ID of the saga instance that is the target
// of the given message.
//
// It returns false if the message should be ignored.
func (m *Mapper) MapMessageToInstance(
	_ context.Context,
	_ saga.Saga,
	_ persistence.Tx,
	env ax.Envelope,
) (saga.InstanceID, bool, error) {
	id, ok := m.Resolver.InstanceIDForMessage(env)
	return id, ok, nil
}

// UpdateMapping notifies the mapper that an instance has been modified,
// allowing it to update it's mapping information, if necessary.
func (m *Mapper) UpdateMapping(
	context.Context,
	saga.Saga,
	persistence.Tx,
	saga.Instance,
) error {
	return nil
}

// DeleteMapping notifies the mapper that an instance has been completed,
// allowing it to remove it's mapping information, if necessary.
func (m *Mapper) DeleteMapping(
	context.Context,
	saga.Saga,
	persistence.Tx,
	saga.Instance,
) error {
	return nil
}
