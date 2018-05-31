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

	res, id, err := h.Mapper.MapMessageToInstance(ctx, h.Saga, tx, env)
	if err != nil {
		return err
	}

	switch res {
	case MapResultFound:
		err = h.handleUpdate(ctx, tx, s, env, id)
	case MapResultNotFound:
		err = h.handleCreate(ctx, tx, s, env)
	case MapResultIgnore:
		return nil
	default:
		panic("unexpected map result type")
	}

	if err != nil {
		return err
	}

	return com.Commit()
}

// handleCreate handles a message that applies to a new saga instance.
// It calls the not-found handler if the message is not a "trigger" message.
func (h *MessageHandler) handleCreate(
	ctx context.Context,
	tx persistence.Tx,
	s ax.Sender,
	env ax.Envelope,
) error {
	triggers, _ := h.Saga.MessageTypes()
	if !triggers.Has(env.Type()) {
		return h.Saga.HandleNotFound(ctx, s, env)
	}

	i, err := h.newInstance(ctx, env)
	if err != nil {
		return err
	}

	w, err := h.Persister.BeginCreate(ctx, h.Saga, tx, s, i)
	if err != nil {
		return err
	}

	return h.handleMessage(ctx, tx, w, env)
}

// handleUpdate handles a message that applies to an existing saga instance.
func (h *MessageHandler) handleUpdate(
	ctx context.Context,
	tx persistence.Tx,
	s ax.Sender,
	env ax.Envelope,
	id InstanceID,
) error {
	w, err := h.Persister.BeginUpdate(ctx, h.Saga, tx, s, id)
	if err != nil {
		return err
	}

	return h.handleMessage(ctx, tx, w, env)
}

// newInstance returns a new saga instance.
func (h *MessageHandler) newInstance(ctx context.Context, env ax.Envelope) (Instance, error) {
	id, err := h.Saga.GenerateInstanceID(ctx, env)
	if err != nil {
		return Instance{}, err
	}

	return Instance{
		InstanceID: id,
		Data:       h.Saga.NewData(),
	}, nil
}

// handleMessage passes the message to the saga for handling, then persists
// any changes to the instance.
func (h *MessageHandler) handleMessage(
	ctx context.Context,
	tx persistence.Tx,
	w UnitOfWork,
	env ax.Envelope,
) error {
	i := w.Instance()
	s := w.Sender()

	if i.Data == nil {
		panic("unit-of-work contains saga instance with nil data")
	}

	if es, ok := h.Saga.(EventedSaga); ok {
		s = &Applier{es, i.Data, s}
	}

	if err := h.Saga.HandleMessage(ctx, s, env, i); err != nil {
		return err
	}

	ok, err := w.Save(ctx)
	if err != nil {
		return err
	}

	if ok {
		if err := h.Mapper.UpdateMapping(ctx, h.Saga, tx, i); err != nil {
			return err
		}
	}

	return nil
}
