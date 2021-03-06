package saga

import (
	"context"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/internal/tracing"
	"github.com/jmalloc/ax/persistence"
	"github.com/opentracing/opentracing-go/log"
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
func (h *MessageHandler) HandleMessage(ctx context.Context, s ax.Sender, mctx ax.MessageContext) error {
	tx, com, err := persistence.GetOrBeginTx(ctx)
	if err != nil {
		return err
	}
	defer com.Rollback()
	ctx = persistence.WithTx(ctx, tx)

	// begin a new unit of work.
	// if ok is false the message is not map to any instance and is ignored.
	w, ok, err := h.begin(ctx, tx, s, mctx.Envelope)
	if err != nil {
		return err
	}

	if !ok {
		h.logEvent(
			ctx,
			"saga_not_mapped",
			"this message does not map to any saga instance",
			nil,
		)

		return nil
	}

	defer w.Close()

	if w.Instance().Revision == 0 {
		if h.isTrigger(mctx.Envelope) {
			h.logEvent(
				ctx,
				"saga_created",
				"this message has triggered a new saga instance",
				w,
			)
		} else {
			h.logEvent(
				ctx,
				"saga_not_found",
				"this message maps to a non-existent saga instance",
				w,
			)

			return h.Saga.HandleNotFound(ctx, s, mctx)
		}
	}

	// otherwise, forward the message to the saga for handling.
	err = h.forward(ctx, w, mctx)
	if err != nil {
		return err
	}

	// check if the instance is complete.
	isComplete, err := h.Saga.IsInstanceComplete(ctx, w.Instance())
	if err != nil {
		return err
	}

	// then persist the changes.
	if isComplete {
		err = h.complete(ctx, tx, w)
	} else {
		err = h.save(ctx, tx, w)
	}

	if err != nil {
		return err
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
func (h *MessageHandler) forward(ctx context.Context, w UnitOfWork, mctx ax.MessageContext) error {
	i := w.Instance()
	s := w.Sender()

	if es, ok := h.Saga.(EventedSaga); ok {
		s = &Applier{es, i.Data, s}
	}

	return h.Saga.HandleMessage(ctx, s, mctx, i)
}

// save persists changes to the saga instance.
func (h *MessageHandler) save(ctx context.Context, tx persistence.Tx, w UnitOfWork) error {
	revBefore := w.Instance().Revision

	ok, err := w.Save(ctx)
	if err != nil {
		return err
	}

	if ok {
		h.logEvent(
			ctx,
			"saga_saved",
			"this message caused a change to the saga instance",
			w,
			log.Uint64("revision_before", uint64(revBefore)),
		)
	} else {
		h.logEvent(
			ctx,
			"saga_not_saved",
			"this message did not cause any change to the saga instance",
			w,
			log.Uint64("revision_before", uint64(revBefore)),
		)
	}

	return h.Mapper.UpdateMapping(ctx, h.Saga, tx, w.Instance())
}

// complete saves a completed saga instance.
func (h *MessageHandler) complete(ctx context.Context, tx persistence.Tx, w UnitOfWork) error {
	revBefore := w.Instance().Revision

	if revBefore != 0 {
		if err := w.SaveAndComplete(ctx); err != nil {
			return err
		}
	}

	h.logEvent(
		ctx,
		"saga_completed",
		"this message completed the saga instance",
		w,
		log.Uint64("revision_before", uint64(revBefore)),
	)

	if revBefore != 0 {
		return h.Mapper.DeleteMapping(ctx, h.Saga, tx, w.Instance())
	}

	return nil
}

// isTrigger returns true if env contains a message type that can trigger a new
// saga instance.
func (h *MessageHandler) isTrigger(env ax.Envelope) bool {
	triggers, _ := h.Saga.MessageTypes()
	return triggers.Has(env.Type())
}

func (h *MessageHandler) logEvent(
	ctx context.Context,
	event, message string,
	w UnitOfWork,
	fields ...log.Field,
) {
	fields = append(
		fields,
		tracing.TypeName("saga", h.Saga),
	)

	if w == nil {
		fields = append(
			fields,
			tracing.TypeName("data", h.Saga.NewData()),
		)
	} else {
		fields = append(
			fields,
			log.String("instance_id", w.Instance().InstanceID.Get()),
			tracing.TypeName("data", w.Instance().Data),
		)

		if w.Instance().Revision > 0 {
			fields = append(
				fields,
				log.Uint64("revision_after", uint64(w.Instance().Revision)),
				log.String("description", w.Instance().Data.InstanceDescription()),
			)
		}
	}

	tracing.LogEvent(
		ctx,
		event,
		message,
		fields...,
	)
}
