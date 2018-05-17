package saga

import (
	"context"

	"github.com/golang/protobuf/proto"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// MessageHandler is an implementation of routing.MessageHandler that handles
// the persistence of saga instances before forwarding the message to a saga.
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

// HandleMessage loads a saga instance, passes env to the saga to be handled, and
// saves the changes to the saga instance.
//
// Changes to the saga are persisted within the outbox transaction if one is
// present in ctx. Otherwise, a new transaction is started.
func (h *MessageHandler) HandleMessage(ctx context.Context, s ax.Sender, env ax.Envelope) error {
	tx, com, err := persistence.GetOrBeginTx(ctx)
	if err != nil {
		return err
	}
	defer com.Rollback()

	// attempt to find an existing saga instance from the message mapping key.
	loadReq := LoadRequest{
		SagaName:    h.Saga.SagaName(),
		MessageType: env.Type(),
		MappingKey:  h.Saga.MapMessage(env.Message),
	}
	res, ok, err := h.Repository.LoadSagaInstance(ctx, tx, loadReq)
	if err != nil {
		return err
	}

	// add our transaction to the context, only for the handle methods.
	handleCtx := persistence.WithTx(ctx, tx)

	saveReq := SaveRequest{
		SagaName: loadReq.SagaName,
	}

	if ok {
		saveReq.InstanceID = res.InstanceID
		saveReq.Instance = res.Instance
		saveReq.CurrentRevision = res.CurrentRevision
	} else {
		triggers, _ := h.Saga.MessageTypes()

		// if no existing instance is found, and this message type does not produce
		// new instances, then the not-found handler is called.
		if !triggers.Has(env.Type()) {
			return h.Saga.HandleNotFound(handleCtx, s, env)
		}

		saveReq.InstanceID, saveReq.Instance = h.Saga.NewInstance(env.Message)
	}

	before := proto.Clone(saveReq.Instance)

	// pass the message to the saga for handling.
	if err := h.Saga.HandleMessage(handleCtx, s, env, saveReq.Instance); err != nil {
		return err
	}

	if !proto.Equal(saveReq.Instance, before) {
		saveReq.MappingTable = buildMappingTable(h.Saga, saveReq.Instance)

		// save the changes to the saga and its mapping table.
		if err := h.Repository.SaveSagaInstance(ctx, tx, saveReq); err != nil {
			return err
		}
	}

	return com.Commit()
}
