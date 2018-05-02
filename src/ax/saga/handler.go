package saga

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// MessageHandler is an implementation of ax.MessageHandler that handles the
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

// HandleMessage loads a saga instance, passes m to the saga to be handled, and
// saves the changes to the saga instance.
//
// Changes to the saga are persisted within the outbox transaction if one is
// present in ctx. Otherwise, a new transaction is started using h.Storage.
func (h *MessageHandler) HandleMessage(ctx ax.MessageContext, m ax.Message) error {
	mt := ax.TypeOf(m)
	mk := h.Saga.MapMessage(m)
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

	hctx := ax.BindContext(
		persistence.WithTx(ctx, tx),
		ctx,
	)

	// if no existing instance is found, and this message type does not produce
	// new instances, then the not-found handler is called.
	if !ok && !h.triggers.Has(mt) {
		return h.Saga.HandleNotFound(hctx, m)
	}

	// pass the message to the saga for handling.
	if err := h.Saga.HandleMessage(hctx, m, si); err != nil {
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
