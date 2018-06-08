package messagestore

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/axmysql/internal/sqlutil"
)

// insertStream inserts a new stream and returns its ID.
//
// s is the name of the stream, n is the initial value of the "next" offset.
func insertStream(
	ctx context.Context,
	tx *sql.Tx,
	s string,
	n uint64,
) (int64, error) {
	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO ax_messagestore_stream SET
			name = ?,
			next = ?`,
		s,
		n,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// incrStreamOffset increments the offset for the given stream by n.
//
// s is the name of the stream. It returns the stream's ID.
// It returns an error if o is not the next free offset.
func incrStreamOffset(
	ctx context.Context,
	tx *sql.Tx,
	s string,
	o uint64,
	n uint64,
) (int64, error) {
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

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if o != next {
		return 0, fmt.Errorf(
			"can not append to stream %s, %d is not the next free offset, expected %d",
			s,
			o,
			next,
		)
	}

	err = sqlutil.ExecSingleRow(
		ctx,
		tx,
		`UPDATE ax_messagestore_stream SET
			next = next + ?
		WHERE stream_id = ?`,
		n,
		id,
	)
	if err != nil {
		return 0, err
	}

	return id, err
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
	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO ax_messagestore_offset SET
			 	next = ?
		ON DUPLICATE KEY UPDATE
			next = next + VALUE(next)`,
		n,
	)
	if err != nil {
		return 0, nil
	}

	ar, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}

	if ar == 1 {
		// if MySQL reports rows affected as 1, that means an insert occurred
		// so we know that the offset was 0
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
			time = ?,
			content_type = ?,
			data = ?`,
		g,
		id,
		o,
		env.Message.MessageDescription(),
		env.MessageID,
		env.CausationID,
		env.CorrelationID,
		env.Time.Format(time.RFC3339Nano),
		contentType,
		data,
	)

	return err
}
