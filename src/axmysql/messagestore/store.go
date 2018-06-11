package messagestore

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/messagestore"
	"github.com/jmalloc/ax/src/ax/persistence"
	mysqlpersistence "github.com/jmalloc/ax/src/axmysql/persistence"
)

// Store is a MySQL-backed implementation of Ax's
// messagestore.GloballyOrderedStore interface.
type Store struct{}

// AppendMessages appends one or more messages to a named stream.
//
// offset is a zero-based index into the stream. An error is returned if
// offset is not the next unused offset in the stream.
func (Store) AppendMessages(
	ctx context.Context,
	ptx persistence.Tx,
	stream string,
	offset uint64,
	envs []ax.Envelope,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

	n := uint64(len(envs))

	var (
		id  int64
		ok  bool
		err error
	)

	if offset == 0 {
		id, ok, err = insertStream(ctx, tx, stream, n)
	} else {
		id, ok, err = incrStreamOffset(ctx, tx, stream, offset, n)
	}

	if err != nil {
		return err
	}

	if !ok {
		// TODO: use OCC error https://github.com/jmalloc/ax/issues/93
		return fmt.Errorf(
			"can not append to stream %s, %d is not the next free offset",
			stream,
			offset,
		)
	}

	g, err := incrGlobalOffset(ctx, tx, n)
	if err != nil {
		return err
	}

	for _, env := range envs {
		if err := insertMessage(
			ctx,
			tx,
			id,
			g,
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
func (Store) OpenStream(
	ctx context.Context,
	ds persistence.DataStore,
	stream string,
	offset uint64,
) (messagestore.Stream, bool, error) {
	db := mysqlpersistence.ExtractDB(ds)

	id, ok, err := lookupStreamID(ctx, db, stream)
	if !ok || err != nil {
		return nil, false, err
	}

	return &Stream{
		Fetcher: &StreamFetcher{
			DB:       db,
			StreamID: id,
		},
		NextOffset: offset,
	}, true, nil
}

// OpenGlobal opens the entire store for reading as a single stream.
//
// The offset may be beyond the end of the stream.
func (Store) OpenGlobal(
	ctx context.Context,
	ds persistence.DataStore,
	offset uint64,
) (messagestore.Stream, error) {
	return &Stream{
		Fetcher: &GlobalFetcher{
			DB: mysqlpersistence.ExtractDB(ds),
		},
		NextOffset: offset,
	}, nil
}

// lookupStreamID returns the ID of the stream named s.
func lookupStreamID(ctx context.Context, db *sql.DB, s string) (int64, bool, error) {
	var id int64

	err := db.QueryRowContext(
		ctx,
		`SELECT
			stream_id
		FROM ax_messagestore_stream
		WHERE name = ?`,
		s,
	).Scan(
		&id,
	)

	if err == sql.ErrNoRows {
		return 0, false, nil
	} else if err != nil {
		return 0, false, err
	}

	return id, true, nil
}
