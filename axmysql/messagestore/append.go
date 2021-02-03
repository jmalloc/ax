package messagestore

import (
	"context"
	"database/sql"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/axmysql/internal/sqlutil"
	"github.com/jmalloc/ax/marshaling"
)

// insertStream inserts a new stream and returns its ID.
//
// s is the name of the stream, n is the initial value of the "next" offset.
// It returns false if there is already a stream by this name.
func insertStream(
	ctx context.Context,
	tx *sql.Tx,
	s string,
	n uint64,
) (int64, bool, error) {
	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO ax_messagestore_stream SET
			name = ?,
			next = ?`,
		s,
		n,
	)
	if err != nil {
		if sqlutil.IsDuplicateEntry(err) {
			return 0, false, nil
		}

		return 0, false, err
	}

	id, err := res.LastInsertId()
	return id, true, err
}

// incrStreamOffset increments the offset for the given stream by n.
//
// s is the name of the stream. It returns the stream's ID.
// It returns false if o is not the next free offset.
func incrStreamOffset(
	ctx context.Context,
	tx *sql.Tx,
	s string,
	o uint64,
	n uint64,
) (int64, bool, error) {
	var (
		id   int64
		next uint64
	)

	err := tx.QueryRowContext(
		ctx,
		`SELECT
			stream_id,
			next
		FROM ax_messagestore_stream
		WHERE name = ?
		FOR UPDATE`, // ensure stream row is locked at this revision
		s,
	).Scan(
		&id,
		&next,
	)

	if err == sql.ErrNoRows {
		return 0, false, nil
	} else if err != nil {
		return 0, false, err
	}

	if o != next {
		return 0, false, nil
	}

	return id, true, sqlutil.ExecSingleRow(
		ctx,
		tx,
		`UPDATE ax_messagestore_stream SET
			next = next + ?
		WHERE stream_id = ?`,
		n,
		id,
	)
}

// incrGlobalOffset increments the global stream offset by n.
//
// It returns the offset before it was incremented.
func incrGlobalOffset(
	ctx context.Context,
	tx *sql.Tx,
	n uint64,
) (uint64, error) {
	// insert or update both cause the row to be locked in tx
	inserted, err := sqlutil.ExecInsertOrUpdate(
		ctx,
		tx,
		`INSERT INTO ax_messagestore_offset SET
			next = ?
		ON DUPLICATE KEY UPDATE
			next = next + VALUE(next)`,
		n,
	)
	if err != nil {
		return 0, nil
	}

	if inserted {
		return 0, nil
	}

	var next uint64
	err = tx.QueryRowContext(
		ctx,
		`SELECT
			next
		FROM ax_messagestore_offset`,
	).Scan(
		&next,
	)

	return next - n, err
}

// insertMessage inserts a message into the store.
//
// id is the stream ID. g and o are the global and stream offsets of the
// message, respectively.
func insertMessage(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
	g uint64,
	o uint64,
	env ax.Envelope,
) error {
	contentType, data, err := ax.MarshalMessage(env.Message)
	if err != nil {
		return err
	}

	// Truncate the message to 255 characters to fit in the schema restriction;
	// see
	// https://github.com/jmalloc/ax/blob/v0.4.0/axmysql/messagestore/schema.sql#L37
	// for details.
	descrTrunc := env.Message.MessageDescription()[:255]

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO ax_messagestore_message SET
			global_offset = ?,
			stream_id = ?,
			stream_offset = ?,
			description = ?,
			message_id = ?,
			causation_id = ?,
			correlation_id = ?,
			created_at = ?,
			send_at = ?,
			content_type = ?,
			data = ?`,
		g,
		id,
		o,
		descrTrunc,
		env.MessageID,
		env.CausationID,
		env.CorrelationID,
		marshaling.MarshalTime(env.CreatedAt),
		marshaling.MarshalTime(env.SendAt),
		contentType,
		data,
	)

	return err
}
