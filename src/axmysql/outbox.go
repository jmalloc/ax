package axmysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// OutboxRepository is an implementation of outbox.Repository that uses SQL
// persistence.
type OutboxRepository struct{}

// LoadOutbox loads the unsent outbound messages that were produced when the
// message identified by id was first delivered.
func (r *OutboxRepository) LoadOutbox(
	ctx context.Context,
	ds persistence.DataStore,
	id ax.MessageID,
) ([]endpoint.OutboundEnvelope, bool, error) {
	db := ds.(*DataStore).DB

	row := db.QueryRowContext(
		ctx,
		`SELECT EXISTS (
			SELECT * FROM outbox WHERE message_id = ?
		)`,
		id,
	)

	var ok bool
	if err := row.Scan(&ok); err != nil {
		return nil, false, err
	}

	if !ok {
		return nil, false, nil
	}

	rows, err := db.QueryContext(
		ctx,
		`SELECT
			message_id,
			correlation_id,
			time,
			content_type,
			body,
			operation,
			destination
		FROM outbox_message
		WHERE causation_id = ?`,
		id,
	)
	if err != nil {
		return nil, false, err
	}

	var envelopes []endpoint.OutboundEnvelope

	for rows.Next() {
		env, err := scanOutboxMessage(rows, id)
		if err != nil {
			return nil, false, err
		}

		envelopes = append(envelopes, env)
	}

	return envelopes, true, nil
}

// SaveOutbox saves a set of unsent outbound messages that were produced
// when the message identified by id was delivered.
func (r *OutboxRepository) SaveOutbox(
	ctx context.Context,
	tx persistence.Tx,
	id ax.MessageID,
	envs []endpoint.OutboundEnvelope,
) error {
	stx := tx.(*Tx).sqlTx

	if _, err := stx.ExecContext(
		ctx,
		`INSERT INTO outbox SET message_id = ?`,
		id,
	); err != nil {
		return err
	}

	for _, env := range envs {
		if err := insertOutboxMessage(ctx, stx, env); err != nil {
			return err
		}
	}

	return nil
}

// MarkAsSent marks a message as sent, removing it from the outbox.
func (r *OutboxRepository) MarkAsSent(
	ctx context.Context,
	tx persistence.Tx,
	env endpoint.OutboundEnvelope,
) error {
	_, err := tx.(*Tx).sqlTx.ExecContext(
		ctx,
		`DELETE FROM outbox_message WHERE message_id = ?`,
		env.MessageID,
	)
	return err
}

func scanOutboxMessage(rows *sql.Rows, causationID ax.MessageID) (endpoint.OutboundEnvelope, error) {
	env := endpoint.OutboundEnvelope{
		Envelope: ax.Envelope{
			CausationID: causationID,
		},
	}

	var (
		ct      string
		body    []byte
		timeStr string
	)

	err := rows.Scan(
		&env.MessageID,
		&env.CorrelationID,
		&timeStr,
		&ct,
		&body,
		&env.Operation,
		&env.DestinationEndpoint,
	)
	if err != nil {
		return env, err
	}

	env.Time, err = time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return env, err
	}

	env.Message, err = ax.UnmarshalMessage(ct, body)

	return env, err
}

func insertOutboxMessage(
	ctx context.Context,
	tx *sql.Tx,
	env endpoint.OutboundEnvelope,
) error {
	ct, body, err := ax.MarshalMessage(env.Message)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO outbox_message SET
			message_id = ?,
			causation_id = ?,
			correlation_id = ?,
			time = ?,
			content_type = ?,
			body = ?,
			operation = ?,
			destination = ?`,
		env.MessageID,
		env.CausationID,
		env.CorrelationID,
		env.Time.Format(time.RFC3339Nano),
		ct,
		body,
		env.Operation,
		env.DestinationEndpoint,
	)

	return err
}

// OutboxSchema is a collection of DDL queries that create the schema
// used by OutboxRepository.
var OutboxSchema = []string{
	`CREATE TABLE IF NOT EXISTS outbox (
	    message_id VARBINARY(255) PRIMARY KEY,
	    time       TIMESTAMP NOT NULL,

	    INDEX (time)
	)`,
	`CREATE TABLE IF NOT EXISTS outbox_message (
	    message_id     VARBINARY(255) NOT NULL PRIMARY KEY,
	    causation_id   VARBINARY(255) NOT NULL, -- outbox.message_id
	    correlation_id VARBINARY(255) NOT NULL,
	    time           VARBINARY(255) NOT NULL,
	    content_type   VARBINARY(255) NOT NULL,
	    body           BLOB NOT NULL,
	    operation      INTEGER NOT NULL,
	    destination    VARBINARY(255) NOT NULL,

	    INDEX (causation_id)
	)`,
}
