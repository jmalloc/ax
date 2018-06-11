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
	Saga      Saga
	Mapper    Mapper
	Persister Persister
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

	// begin a new unit of work.
	// if ok is false the message is not map to any instance and is ignored.
	w, ok, err := h.begin(ctx, tx, s, env)
	if !ok || err != nil {
		return err
	}
	defer w.Close()

	// call the not-found handler if the instance is new but the message is not a
	// trigger.
	if w.Instance().Revision == 0 && !h.isTrigger(env) {
		return h.Saga.HandleNotFound(ctx, s, env)
	}

	// otherwise, forward the message to the saga for handling.
	if err := h.forward(ctx, w, env); err != nil {
		return err
	}

	// check if the instance is complete.
	isComplete, err := h.isComplete(ctx, w.Instance())
	if err != nil {
		return err
	}

	// then persist the changes.
	if isComplete {
		if err := h.complete(ctx, tx, w); err != nil {
			return err
		}
	} else {
		if err := h.save(ctx, tx, w); err != nil {
			return err
		}
	}

	return com.Commit()
}

// begin starts a new unit-of-work.
//
// It returns false if the message is not mapped to any instance and hence
// should be ignored.
func (h *MessageHandler) begin(
	ctx context.Context,
	tx persistence.Tx,
	s ax.Sender,
	env ax.Envelope,
) (UnitOfWork, bool, error) {
	id, ok, err := h.Mapper.MapMessageToInstance(ctx, h.Saga, tx, env)
	if !ok || err != nil {
		return nil, false, err
	}

	w, err := h.Persister.BeginUnitOfWork(ctx, h.Saga, tx, s, id)
	return w, true, err
}

// forward passes the message to the saga to be handled.
func (h *MessageHandler) forward(ctx context.Context, w UnitOfWork, env ax.Envelope) error {
	i := w.Instance()
	s := w.Sender()

	if es, ok := h.Saga.(EventedSaga); ok {
		s = &Applier{es, i.Data, s}
	}

	return h.Saga.HandleMessage(ctx, s, env, i)
}

// save persists changes to the saga instance.
func (h *MessageHandler) save(ctx context.Context, tx persistence.Tx, w UnitOfWork) error {
	ok, err := w.Save(ctx)
	if !ok || err != nil {
		return err
	}

	return h.Mapper.UpdateMapping(ctx, h.Saga, tx, w.Instance())
}

// complete saves a completed saga instance.
func (h *MessageHandler) complete(ctx context.Context, tx persistence.Tx, w UnitOfWork) error {
	if err := w.SaveAndComplete(ctx); err != nil {
		return err
	}

	return h.Mapper.DeleteMapping(ctx, h.Saga, tx, w.Instance())
}

// isTrigger returns true if env contains a message type that can trigger a new
// saga instance.
func (h *MessageHandler) isTrigger(env ax.Envelope) bool {
	triggers, _ := h.Saga.MessageTypes()
	return triggers.Has(env.Type())
}

// isComplete returns true if h.Saga is a completable saga and i is complete.
func (h *MessageHandler) isComplete(ctx context.Context, i Instance) (bool, error) {
	if cs, ok := h.Saga.(CompletableSaga); ok {
		return cs.IsComplete(ctx, i)
	}

	return false, nil
}
