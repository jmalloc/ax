package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// MessageHandler is an implementation of bus.MessageHandler that handles the
// persistence of saga instances before forwarding the message to a saga.
type MessageHandler struct {
	Repository Repository
	Saga       Saga

	triggers ax.MessageTypeSet
}

// MessageTypes returns the set of messages that the handler can handle.
//
// For sagas, this is the union of the message types that trigger new instances
// and the message types that are routed to existing instances.
func (h *MessageHandler) MessageTypes() ax.MessageTypeSet {
	triggers, others := h.Saga.MessageTypes()
	h.triggers = triggers

	return triggers.Union(others)
}

// HandleMessage loads a saga instance, passes env to the saga to be handled, and
// saves the changes to the saga instance.
//
// Changes to the saga are persisted within the outbox transaction if one is
// present in ctx. Otherwise, a new transaction is started.
func (h *MessageHandler) HandleMessage(ctx context.Context, s ax.Sender, env ax.Envelope) error {
	mt := env.Type()
	mk := h.Saga.MapMessage(env.Message)
	si := h.Saga.InitialState()

	// attempt to find an existing saga instance from the message mapping key.
	ok, err := h.Repository.LoadSagaInstance(ctx, mt, mk, si)
	if err != nil {
		return err
	}

	tx, com, err := persistence.GetOrBeginTx(ctx)
	if err != nil {
		return err
	}
	defer com.Rollback()

	hctx := persistence.WithTx(ctx, tx)

	// if no existing instance is found, and this message type does not produce
	// new instances, then the not-found handler is called.
	if !ok && !h.triggers.Has(mt) {
		return h.Saga.HandleNotFound(hctx, s, env)
	}

	// pass the message to the saga for handling.
	if err := h.Saga.HandleMessage(hctx, s, env, si); err != nil {
		return err
	}

	// save the changes to the saga and its mapping table.
	if err := h.Repository.SaveSagaInstance(
		ctx,
		tx,
		si,
		buildMappingTable(h.Saga, si),
	); err != nil {
		return err
	}

	return com.Commit()
}
