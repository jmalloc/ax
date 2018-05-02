package dialect

import (
	"context"
	"database/sql"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/ax/marshaling"
)

func (mysql) SelectOutboxExists(
	ctx context.Context,
	db *sql.DB,
	id ax.MessageID,
) (bool, error) {
	row := db.QueryRowContext(
		ctx,
		"SELECT EXISTS (SELECT * FROM `outbox` WHERE `message_id` = ?)",
		id,
	)

	var ok bool
	err := row.Scan(&ok)
	return ok, err
}

func (mysql) SelectOutboxEnvelopes(
	ctx context.Context,
	db *sql.DB,
	id ax.MessageID,
) ([]bus.OutboundEnvelope, error) {
	rows, err := db.QueryContext(
		ctx,
		"SELECT "+
			"`message_id`, "+
			"`correlation_id`, "+
			"`time`, "+
			"`content_type`, "+
			"`body`, "+
			"`operation`, "+
			"`destination` "+
			"FROM `outbox_envelope` WHERE `causation_id` = ?",
		id,
	)
	if err != nil {
		return nil, err
	}

	return scanOutboxMessages(rows, id)
}

func (mysql) InsertOutbox(
	ctx context.Context,
	tx *sql.Tx,
	id ax.MessageID,
) error {
	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO `outbox` SET `message_id` = ?",
		id,
	)

	return err
}

func (mysql) InsertOutboxEnvelope(
	ctx context.Context,
	tx *sql.Tx,
	m bus.OutboundEnvelope,
) error {
	ct, body, err := marshaling.MarshalMessage(m.Message)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO `outbox_envelope` SET "+
			"`message_id` = ?, "+
			"`causation_id` = ?, "+
			"`correlation_id` = ?, "+
			"`time` = ?, "+
			"`content_type` = ?, "+
			"`body` = ?, "+
			"`operation` = ?, "+
			"`destination` = ? ",
		m.MessageID,
		m.CausationID,
		m.CorrelationID,
		m.Time,
		ct,
		body,
		m.Operation,
		m.DestinationEndpoint,
	)

	return err
}

func (mysql) DeleteOutboxEnvelope(
	ctx context.Context,
	tx *sql.Tx,
	id ax.MessageID,
) error {
	_, err := tx.ExecContext(
		ctx,
		"DELETE FROM `outbox_envelope` WHERE `message_id` = ?",
		id,
	)

	return err
}
