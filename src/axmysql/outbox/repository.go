package outbox

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/persistence"
	mysqlpersistence "github.com/jmalloc/ax/src/axmysql/persistence"
)

// Repository is a MySQL-backed implementation of Ax's outbox.Repository
// interface.
type Repository struct{}

// LoadOutbox loads the unsent outbound messages that were produced when the
// message identified by id was first delivered.
func (Repository) LoadOutbox(
	ctx context.Context,
	ds persistence.DataStore,
	id ax.MessageID,
) ([]endpoint.OutboundEnvelope, bool, error) {
	db := mysqlpersistence.ExtractDB(ds)

	row := db.QueryRowContext(
		ctx,
		`SELECT EXISTS (
			SELECT * FROM outbox WHERE causation_id = ?
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
	defer rows.Close()

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
func (Repository) SaveOutbox(
	ctx context.Context,
	ptx persistence.Tx,
	id ax.MessageID,
	envs []endpoint.OutboundEnvelope,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO outbox SET causation_id = ?`,
		id,
	); err != nil {
		return err
	}

	for _, env := range envs {
		if err := insertOutboxMessage(ctx, tx, env); err != nil {
			return err
		}
	}

	return nil
}

// MarkAsSent marks a message as sent, removing it from the outbox.
func (Repository) MarkAsSent(
	ctx context.Context,
	ptx persistence.Tx,
	env endpoint.OutboundEnvelope,
) error {
	_, err := mysqlpersistence.ExtractTx(ptx).ExecContext(
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
