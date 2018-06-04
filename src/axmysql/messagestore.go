package axmysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/messagestore"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// MessageStore is an implementation of messagestore.Store that uses SQL
// persistence.
type MessageStore struct{}

// AppendMessages appends one or more messages to a named stream.
//
// offset is a zero-based index into the stream. An error is returned if
// offset is not the next unused offset in the stream.
func (s MessageStore) AppendMessages(
	ctx context.Context,
	ptx persistence.Tx,
	stream string,
	offset uint64,
	envs []ax.Envelope,
) error {
	tx := sqlTx(ptx)

	var (
		id  int64
		err error
	)

	n := uint64(len(envs))

	if offset == 0 {
		id, err = s.insertStream(ctx, tx, stream, n)
	} else {
		id, err = s.incrStreamOffset(ctx, tx, stream, offset, n)
	}

	if err != nil {
		return err
	}

	g, err := s.incrGlobalOffset(ctx, tx, n)
	if err != nil {
		return err
	}

	for _, env := range envs {
		if err := s.insertStreamMessage(
			ctx,
			tx,
			g,
			id,
			offset,
			env,
		); err != nil {
			return err
		}

		g++
		offset++
	}

	return nil
}

// OpenStream opens a stream of messages for reading from a specific offset.
//
// The offset may be past the end of the stream. It returns false if the stream
// does not exist.
func (MessageStore) OpenStream(
	ctx context.Context,
	pds persistence.DataStore,
	stream string,
	offset uint64,
) (messagestore.Stream, bool, error) {
	db := pds.(*DataStore).DB

	var id int64

	if err := db.QueryRowContext(
		ctx,
		`SELECT
			stream_id
		FROM messagestore_stream
		WHERE name = ?`,
		stream,
	).Scan(
		&id,
	); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		return nil, false, err
	}

	return &MessageStream{
		db:     db,
		id:     id,
		offset: offset,
	}, true, nil
}

// insertStream inserts a new stream and returns its ID.
//
// s is the name of the stream, n is the number of events being appended.
func (MessageStore) insertStream(
	ctx context.Context,
	tx *sql.Tx,
	s string,
	n uint64,
) (int64, error) {
	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO messagestore_stream SET
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
// s is the name of the stream. It returns an error if o is not the next free
// offset.
func (MessageStore) incrStreamOffset(
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
		FROM messagestore_stream
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

	_, err = tx.ExecContext(
		ctx,
		`UPDATE messagestore_stream SET
			next = next + ?
		WHERE stream_id = ?`,
		n,
		id,
	)

	return id, err
}

// incrGlobalOffset increments the global stream offset by n.
//
// It returns the offset before it was incremented.
func (MessageStore) incrGlobalOffset(
	ctx context.Context,
	tx *sql.Tx,
	n uint64,
) (uint64, error) {
	// insert or update both cause the row to be locked in tx
	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO messagestore_offset SET
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
		FROM messagestore_offset`,
	).Scan(
		&next,
	)

	return next - n, err
}

// insertStreamMessage appends a message to a stream.
//
// g is the message's global offset, id is the stream's id and o is the
// message's offset within the stream.
func (MessageStore) insertStreamMessage(
	ctx context.Context,
	tx *sql.Tx,
	g uint64,
	id int64,
	o uint64,
	env ax.Envelope,
) error {
	contentType, data, err := ax.MarshalMessage(env.Message)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO messagestore_message SET
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

// streamFetchLimit is the number of messages to fetch in each select query on
// a message stream.
const streamFetchLimit = 100

// MessageStream is a stream of messages stored in an SQL MessageStore.
type MessageStream struct {
	db     *sql.DB
	id     int64
	offset uint64
	rows   *sql.Rows
}

// TryNext advances the stream to the next message.
//
// It returns false if there are no more messages in the stream.
func (s *MessageStream) TryNext(ctx context.Context) (bool, error) {
	if s.rows != nil {
		if s.advance() {
			return true, nil
		}
	}

	if err := s.fetchRows(ctx); err != nil {
		return false, err
	}

	return s.advance(), nil
}

// Get returns the message at the current offset in the stream.
func (s *MessageStream) Get(ctx context.Context) (ax.Envelope, error) {
	if s.rows == nil {
		panic("Next() must be called before Get()")
	}

	var (
		env         ax.Envelope
		contentType string
		data        []byte
		timeStr     string
	)

	err := s.rows.Scan(
		&env.MessageID,
		&env.CausationID,
		&env.CorrelationID,
		&timeStr,
		&contentType,
		&data,
	)
	if err != nil {
		return ax.Envelope{}, err
	}

	env.Time, err = time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return ax.Envelope{}, err
	}

	env.Message, err = ax.UnmarshalMessage(contentType, data)

	return env, err
}

// Offset returns the offset of the message returned by Get().
func (s *MessageStream) Offset() (uint64, error) {
	return s.offset, nil
}

// Close closes the stream.
func (s *MessageStream) Close() error {
	return s.replaceRows(nil)
}

// fetchRows selects the next batch of messages from the stream.
func (s *MessageStream) fetchRows(ctx context.Context) error {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT
			message_id,
			causation_id,
			correlation_id,
			time,
			content_type,
			data
		FROM messagestore_message
		WHERE stream_id = ?
		AND stream_offset >= ?
		ORDER BY stream_offset
		LIMIT ?`,
		s.id,
		s.offset,
		streamFetchLimit,
	)
	if err != nil {
		return err
	}

	return s.replaceRows(rows)
}

// replaceRows replaces s.rows with r, closing the existing s.rows value if it
// is not nil.
func (s *MessageStream) replaceRows(r *sql.Rows) error {
	if s.rows != nil {
		prev := s.rows
		s.rows = nil

		if err := prev.Close(); err != nil {
			return err
		}
	}

	s.rows = r
	return nil
}

// advance moves to the next row in s.rows.
func (s *MessageStream) advance() bool {
	if s.rows.Next() {
		s.offset++
		return true
	}

	return false
}
