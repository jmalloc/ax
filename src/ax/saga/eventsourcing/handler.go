package eventsourcing

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// MessageHandler is an implementation of routing.MessageHandler that loads a
// saga instance, forwards the message to the saga, then perists any changes
// to the instance.
type MessageHandler struct {
	Saga       saga.Saga
	Mapper     saga.Mapper
	Repository Repository
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

	i, ok, err := h.findInstance(ctx, tx, env)
	if err != nil {
		return err
	}

	if ok {
		sender := &Sender{
			Next: s,
			Data: i.Data.(Data),
		}

		if err = h.Saga.HandleMessage(ctx, sender, env, i); err != nil {
			return err
		}

		if err := h.saveInstance(ctx, tx, i, sender.Events); err != nil {
			return err
		}
	} else {
		if err := h.Saga.HandleNotFound(ctx, s, env); err != nil {
			return err
		}
	}

	return com.Commit()
}

// findInstance returns the saga instance that the given message is routed to,
// creating a new instance if necessary.
func (h *MessageHandler) findInstance(
	ctx context.Context,
	tx persistence.Tx,
	env ax.Envelope,
) (saga.Instance, bool, error) {
	k, err := h.Saga.MappingKeyForMessage(ctx, env)
	if err != nil {
		return saga.Instance{}, false, err
	}

	sn := h.Saga.SagaName()
	id, ok, err := h.Mapper.FindByKey(ctx, tx, sn, k)
	if err != nil {
		return saga.Instance{}, false, err
	}

	if ok {
		zv, err := h.Saga.InitialState(ctx)
		if err != nil {
			return saga.Instance{}, false, err
		}

		i, err := h.Repository.LoadSagaInstance(ctx, tx, id, zv.(Data))
		return i, true, err
	}

	triggers, _ := h.Saga.MessageTypes()
	if triggers.Has(env.Type()) {
		i, err := h.newInstance(ctx, env)
		return i, true, err
	}

	return saga.Instance{}, false, nil
}

func (h *MessageHandler) newInstance(ctx context.Context, env ax.Envelope) (saga.Instance, error) {
	id, err := h.Saga.GenerateInstanceID(ctx, env)
	if err != nil {
		return saga.Instance{}, err
	}

	data, err := h.Saga.InitialState(ctx)
	if err != nil {
		return saga.Instance{}, err
	}

	return saga.Instance{
		InstanceID: id,
		Data:       data,
	}, nil
}

// saveInstance persists updates to an instance's data and mapping key set.
func (h *MessageHandler) saveInstance(
	ctx context.Context,
	tx persistence.Tx,
	i saga.Instance,
	ev []ax.Event,
) error {
	if len(ev) == 0 {
		return nil
	}

	if err := h.Repository.SaveSagaInstance(ctx, tx, i, ev); err != nil {
		return err
	}

	ks, err := h.Saga.MappingKeysForInstance(ctx, i)
	if err != nil {
		return err
	}

	sn := h.Saga.SagaName()
	return h.Mapper.SaveKeys(ctx, tx, sn, i.InstanceID, ks)
}
