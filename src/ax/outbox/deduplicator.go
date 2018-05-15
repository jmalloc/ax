package outbox

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// Deduplicator is an inbound pipeline stage that provides message idempotency
// using the "outbox" pattern.
//
// See http://gistlabs.com/2014/05/the-outbox/
type Deduplicator struct {
	Repository Repository
	Next       endpoint.InboundPipeline
}

// Initialize calls d.Next.Initialize()
func (d *Deduplicator) Initialize(ctx context.Context, t endpoint.Transport) error {
	return d.Next.Initialize(ctx, t)
}

// Accept passes env to the next pipeline stage only if it has not been
// delivered previously.
//
// If it has been delivered previously, the messages that were produced the
// first time are sent using s.
func (d *Deduplicator) Accept(ctx context.Context, s endpoint.MessageSink, env endpoint.InboundEnvelope) error {
	ds, ok := persistence.GetDataStore(ctx)
	if !ok {
		return errors.New("no data store is available in ctx")
	}

	messages, ok, err := d.Repository.LoadOutbox(
		ctx,
		ds,
		env.MessageID,
	)
	if err != nil {
		return err
	}

	if !ok {
		messages, err = d.forward(ctx, env)
		if err != nil {
			return err
		}
	}

	for _, o := range messages {
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

	if err := d.Repository.SaveOutbox(
		ctx,
		tx,
		env.MessageID,
		s.Envelopes,
	); err != nil {
		return nil, err
	}

	return s.Envelopes, com.Commit()
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
