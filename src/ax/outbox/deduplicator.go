package outbox

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// Deduplicator is an inbound pipeline stage that provides message idempotency
// using the "outbox" pattern.
//
// See http://gistlabs.com/2014/05/the-outbox/
type Deduplicator struct {
	Repository Repository
	Next       bus.InboundPipeline
}

// Initialize calls d.Next.Initialize()
func (d *Deduplicator) Initialize(ctx context.Context, t bus.Transport) error {
	return d.Next.Initialize(ctx, t)
}

// DeliverMessage passes m to the next pipeline stage only if it has not been
// delivered previously.
//
// If it has been delivered previously, the messages that were produced the
// first time are sent using s.
func (d *Deduplicator) DeliverMessage(
	ctx context.Context,
	s bus.MessageSender,
	m bus.InboundEnvelope,
) error {
	ds, ok := persistence.GetDataStore(ctx)
	if !ok {
		return errors.New("no data store is available in ctx")
	}

	messages, ok, err := d.Repository.LoadOutbox(
		ctx,
		ds,
		m.Envelope.MessageID,
	)
	if err != nil {
		return err
	}

	if !ok {
		messages, err = d.deliver(ctx, m)
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

// deliver passes m to the next pipeline stage and captures and persists the
// messages it produces.
func (d *Deduplicator) deliver(ctx context.Context, m bus.InboundEnvelope) ([]bus.OutboundEnvelope, error) {
	tx, com, err := persistence.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer com.Rollback()

	var s bus.MessageBuffer

	if err := d.Next.DeliverMessage(
		persistence.WithTx(ctx, tx),
		&s,
		m,
	); err != nil {
		return nil, err
	}

	if err := d.Repository.SaveOutbox(
		ctx,
		tx,
		m.Envelope.MessageID,
		s.Messages,
	); err != nil {
		return nil, err
	}

	return s.Messages, com.Commit()
}

// send uses s to send a message that was previously persisted before marking it
// as sent.
func (d *Deduplicator) send(
	ctx context.Context,
	s bus.MessageSender,
	m bus.OutboundEnvelope,
) error {
	if err := s.SendMessage(ctx, m); err != nil {
		return err
	}

	tx, com, err := persistence.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer com.Rollback()

	if err := d.Repository.MarkAsSent(ctx, tx, m); err != nil {
		return err
	}

	return com.Commit()
}
