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
	k, ok, err := h.Saga.MappingKeyForMessage(ctx, env)
	if !ok || err != nil {
		return err
	}

	tx, com, err := persistence.GetOrBeginTx(ctx)
	if err != nil {
		return err
	}
	defer com.Rollback()
	ctx = persistence.WithTx(ctx, tx)

	w, ok, err := h.startUnitOfWork(ctx, tx, k, s, env)
	if err != nil {
		return err
	}

	if ok {
		err = h.handleMessage(ctx, tx, w, env)
	} else {
		err = h.Saga.HandleNotFound(ctx, s, env)
	}

	if err != nil {
		return err
	}

	return com.Commit()
}

// startUnitOfWork begins a new unit of work, creating a new saga instance if
// necessary.
func (h *MessageHandler) startUnitOfWork(
	ctx context.Context,
	tx persistence.Tx,
	k string,
	s ax.Sender,
	env ax.Envelope,
) (UnitOfWork, bool, error) {
	id, ok, err := h.Mapper.FindByKey(
		ctx,
		tx,
		h.Saga.SagaName(),
		k,
	)
	if err != nil {
		return nil, false, err
	}

	if ok {
		var w UnitOfWork
		w, err = h.Persister.BeginUpdate(ctx, h.Saga, tx, s, id)
		return w, true, err
	}

	triggers, _ := h.Saga.MessageTypes()
	if !triggers.Has(env.Type()) {
		return nil, false, nil
	}

	i, err := h.newInstance(ctx, env)
	if err != nil {
		return nil, false, err
	}

	w, err := h.Persister.BeginCreate(ctx, h.Saga, tx, s, i)
	return w, true, err
}

// saveKeySet rebuilds and persists the mapping key set for the given instance.
func (h *MessageHandler) saveKeySet(
	ctx context.Context,
	tx persistence.Tx,
	i Instance,
) error {
	ks, err := h.Saga.MappingKeysForInstance(ctx, i)
	if err != nil {
		return err
	}

	return h.Mapper.SaveKeys(
		ctx,
		tx,
		h.Saga.SagaName(),
		i.InstanceID,
		ks,
	)
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
		if err := h.saveKeySet(ctx, tx, i); err != nil {
			return err
		}
	}

	return nil
}
