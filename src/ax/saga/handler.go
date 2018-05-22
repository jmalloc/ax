package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// MessageHandler is an implementation of routing.MessageHandler that loads a
// saga instance, forwards the message to the saga, then perists any changes
// to the instance.
type MessageHandler struct {
	Repository Repository
	Saga       Saga
}

// MessageTypes returns the set of messages that the handler can handle.
//
// For sagas, this is the union of the message types that trigger new instances
// and the message types that are routed to existing instances.
func (h *MessageHandler) MessageTypes() ax.MessageTypeSet {
	triggers, others := h.Saga.MessageTypes()
	return triggers.Union(others)
}

// HandleMessage loads a saga instance, passes env to the saga to be handled,
// and saves the changes to the saga instance.
//
// Changes to the saga are persisted within the existing transaction in ctx, if
// present.
func (h *MessageHandler) HandleMessage(ctx context.Context, s ax.Sender, env ax.Envelope) error {
	tx, com, err := persistence.GetOrBeginTx(ctx)
	if err != nil {
		return err
	}
	defer com.Rollback()

	ctx = persistence.WithTx(ctx, tx)

	// acquire the key used to query the repository.
	k, err := h.Saga.MappingKeyForMessage(ctx, env)
	if err != nil {
		return err
	}

	// attempt to find an existing saga instance.
	sn := h.Saga.SagaName()
	i, ok, err := h.Repository.LoadSagaInstance(ctx, tx, sn, k)
	if err != nil {
		return err
	}

	if !ok {
		triggers, _ := h.Saga.MessageTypes()

		// if no existing instance is found, and this message type does not trigger
		// new instances, then the not-found handler is called.
		if !triggers.Has(env.Type()) {
			return h.Saga.HandleNotFound(ctx, s, env)
		}

		// otherwise, create a new saga instance.
		i.InstanceID, i.Data, err = h.Saga.NewInstance(ctx, env)
		if err != nil {
			return err
		}
	}

	// pass the message to the saga for handling.
	err = h.Saga.HandleMessage(ctx, s, env, i)
	if err != nil {
		return err
	}

	// rebuild the instance's mapping table.
	ks, err := h.Saga.MappingKeysForInstance(ctx, i)
	if err != nil {
		return err
	}

	// save the changes to the saga and its mapping table.
	if err := h.Repository.SaveSagaInstance(ctx, tx, sn, i, ks); err != nil {
		return err
	}

	return com.Commit()
}
