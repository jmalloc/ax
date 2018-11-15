package outbox

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/internal/tracing"
	"github.com/opentracing/opentracing-go/log"
)

// Deduplicator is an inbound pipeline stage that provides message idempotency
// using the "outbox" pattern.
//
// See http://gistlabs.com/2014/05/the-outbox/
type Deduplicator struct {
	Repository Repository
	Next       endpoint.InboundPipeline
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (d *Deduplicator) Initialize(ctx context.Context, ep *endpoint.Endpoint) error {
	return d.Next.Initialize(ctx, ep)
}

// Accept passes env to the next pipeline stage only if it has not been
// processed previously.
//
// If it has been processed previously, the messages that were produced the
// first time are sent using s.
func (d *Deduplicator) Accept(ctx context.Context, s endpoint.MessageSink, env endpoint.InboundEnvelope) error {
	ds, ok := persistence.GetDataStore(ctx)
	if !ok {
		return errors.New("no data store is available in ctx")
	}

	envs, ok, err := d.Repository.LoadOutbox(
		ctx,
		ds,
		env.MessageID,
	)
	if err != nil {
		return err
	}

	if ok {
		tracing.LogEvent(
			ctx,
			"outbox-loaded",
			"found an outbox, the message has already been processed",
			log.Int("pending-outbound-messages", len(envs)),
			tracing.TypeName("pipeline-stage", d),
		)
	} else {
		tracing.LogEvent(
			ctx,
			"outbox-not-found",
			"no outbox found, forwarding message to the next pipeline stage",
			tracing.TypeName("pipeline-stage", d),
		)

		envs, err = d.forward(ctx, env)
		if err != nil {
			return err
		}

		tracing.LogEvent(
			ctx,
			"outbox-saved",
			"message processed successfully, saving outbox",
			log.Int("pending-outbound-messages", len(envs)),
			tracing.TypeName("pipeline-stage", d),
		)
	}

	for _, o := range envs {
		if err := d.send(ctx, s, o); err != nil {
			return err
		}
	}

	return nil
}

// forward passes env to the next pipeline stage and persists the messages it produces to the outbox.
// The messages are also returned to be sent via the transport immediately.
func (d *Deduplicator) forward(ctx context.Context, env endpoint.InboundEnvelope) ([]endpoint.OutboundEnvelope, error) {
	tx, com, err := persistence.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer com.Rollback()

	var s endpoint.BufferedSink

	if err := d.Next.Accept(
		persistence.WithTx(ctx, tx),
		&s,
		env,
	); err != nil {
		return nil, err
	}

	envs := s.TakeEnvelopes()

	if err := d.Repository.SaveOutbox(
		ctx,
		tx,
		env.MessageID,
		envs,
	); err != nil {
		return nil, err
	}

	return envs, com.Commit()
}

// send uses s to send a message that was previously persisted before marking it
// as sent.
func (d *Deduplicator) send(
	ctx context.Context,
	s endpoint.MessageSink,
	env endpoint.OutboundEnvelope,
) error {
	if err := s.Accept(ctx, env); err != nil {
		return err
	}

	tx, com, err := persistence.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer com.Rollback()

	if err := d.Repository.MarkAsSent(ctx, tx, env); err != nil {
		return err
	}

	return com.Commit()
}
