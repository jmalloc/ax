package envelopestore

import (
	"context"
	"database/sql"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/marshaling"
)

// Insert adds a message to the store.
func Insert(
	ctx context.Context,
	tx *sql.Tx,
	table string,
	env endpoint.OutboundEnvelope,
) error {
	ct, data, err := ax.MarshalMessage(env.Message)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO `+table+` SET
			message_id = ?,
			causation_id = ?,
			correlation_id = ?,
			created_at = ?,
			send_at = ?,
			content_type = ?,
			data = ?,
			operation = ?,
			destination = ?`,
		env.MessageID,
		env.CausationID,
		env.CorrelationID,
		marshaling.MarshalTime(env.CreatedAt),
		marshaling.MarshalTime(env.SendAt),
		ct,
		data,
		env.Operation,
		env.DestinationEndpoint,
	)

	return err
}

// Delete removes a message from the store.
func Delete(
	ctx context.Context,
	tx *sql.Tx,
	table string,
	env endpoint.OutboundEnvelope,
) error {
	_, err := tx.ExecContext(
		ctx,
		`DELETE FROM `+table+` WHERE message_id = ?`,
		env.MessageID,
	)

	return err
}
