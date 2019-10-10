package keyset

import (
	"context"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/persistence"
	"github.com/jmalloc/ax/saga"
)

// Mapper is an implementation of saga.Mapper that maps messages to saga
// instances using disjoint "key sets".
//
// This is a flexible mapping strategy that allows the saga precise control over
// which instance to load on a per-message basis.
//
// The saga must implement the keyset.Saga interface to use key set mapping.
type Mapper struct {
	Repository Repository
	Resolver   Resolver
}

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
	k, ok, err := m.Resolver.MappingKeyForMessage(ctx, env)
	if !ok || err != nil {
		return saga.InstanceID{}, false, err
	}

	id, ok, err := m.Repository.FindByKey(ctx, tx, sg.PersistenceKey(), k)
	if ok || err != nil {
		return id, true, err
	}

	id, err = m.Resolver.GenerateInstanceID(ctx, env)
	return id, true, err
}

// UpdateMapping notifies the mapper that an instance has been modified,
// allowing it to update it's mapping information, if necessary.
func (m *Mapper) UpdateMapping(
	ctx context.Context,
	sg saga.Saga,
	tx persistence.Tx,
	i saga.Instance,
) error {
	ks, err := m.Resolver.MappingKeysForInstance(ctx, i)
	if err != nil {
		return err
	}

	ks, err = Validate(ks)
	if err != nil {
		return err
	}

	return m.Repository.SaveKeys(ctx, tx, sg.PersistenceKey(), ks, i.InstanceID)
}

// DeleteMapping notifies the mapper that an instance has been completed,
// allowing it to remove it's mapping information, if necessary.
func (m *Mapper) DeleteMapping(
	ctx context.Context,
	sg saga.Saga,
	tx persistence.Tx,
	i saga.Instance,
) error {
	return m.Repository.DeleteKeys(ctx, tx, sg.PersistenceKey(), i.InstanceID)
}
