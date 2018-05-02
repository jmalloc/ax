package axsql

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// OutboxRepository is an implementation of outbox.Repository that uses SQL
// persistence.
type OutboxRepository struct{}

// LoadOutbox loads the unsent outbound messages that were produced when the
// message identified by id was first delivered.
//
// ok is false if the message has not yet been successfully delivered.
func (r *OutboxRepository) LoadOutbox(
	ctx context.Context,
	id ax.MessageID,
) (m []bus.OutboundEnvelope, ok bool, err error) {
	v, ok := persistence.GetDataStore(ctx)
	if !ok {
		err = errors.New("no data store is available in ctx")
		return
	}

	ds, ok := v.(*DataStore)
	if !ok {
		err = errors.New("data store in ctx is not an axsql.DataStore")
		return
	}

	ok, err = ds.Dialect.SelectOutboxExists(ctx, ds.DB, id)
	if err != nil || !ok {
		return
	}

	m, err = ds.Dialect.SelectOutboxEnvelopes(ctx, ds.DB, id)
	return
}

// SaveOutbox saves a set of unsent outbound messages that were produced
// when the message identified by id was delivered.
func (r *OutboxRepository) SaveOutbox(
	ctx context.Context,
	tx persistence.Tx,
	id ax.MessageID,
	m []bus.OutboundEnvelope,
) error {
	t, ok := tx.(*Tx)
	if !ok {
		return errors.New("transaction in ctx is not an *axsql.Tx")
	}

	if err := t.ds.Dialect.InsertOutbox(ctx, t.tx, id); err != nil {
		return err
	}

	for _, env := range m {
		if err := t.ds.Dialect.InsertOutboxEnvelope(ctx, t.tx, env); err != nil {
			return err
		}
	}

	return nil
}

// MarkAsSent marks a message as sent, removing it from the outbox.
func (r *OutboxRepository) MarkAsSent(
	ctx context.Context,
	tx persistence.Tx,
	m bus.OutboundEnvelope,
) error {
	t, ok := tx.(*Tx)
	if !ok {
		return errors.New("transaction in ctx is not an *axsql.Tx")
	}

	return t.ds.Dialect.DeleteOutboxEnvelope(ctx, t.tx, m.MessageID)
}
