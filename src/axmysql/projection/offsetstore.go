package projection

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axmysql/internal/sqlutil"
	mysqlpersistence "github.com/jmalloc/ax/src/axmysql/persistence"
)

// OffsetStore is a MySQL-backed implementation of Ax's projection.OffsetStore
// interface.
type OffsetStore struct{}

// LoadOffset returns the offset at which a consumer should resume
// reading from the stream.
//
// pk is the projector's persistence key.
func (OffsetStore) LoadOffset(
	ctx context.Context,
	ds persistence.DataStore,
	pk string,
) (uint64, error) {
	db := mysqlpersistence.ExtractDB(ds)

	var offset uint64

	err := db.QueryRowContext(
		ctx,
		`SELECT
			next_offset
		FROM ax_projection_offset
		WHERE persistence_key = ?`,
		pk,
	).Scan(
		&offset,
	)

	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return offset, nil
}

// IncrementOffset increments the offset at which a consumer should resume
// reading from the stream by one.
//
// pk is the projector's persitence key. c is the offset that is currently
// stored, as returned by LoadOffset(). If c is not the offset that is
// currently stored, the increment fails and a non-nil error is returned.
func (s OffsetStore) IncrementOffset(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	c uint64,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

	var (
		ok  bool
		err error
	)

	if c == 0 {
		ok, err = s.insertOffset(ctx, tx, pk)
	} else {
		ok, err = s.updateOffset(ctx, tx, pk, c)
	}

	if ok || err != nil {
		return err
	}

	// TODO: use OCC error https://github.com/jmalloc/ax/issues/93
	return fmt.Errorf(
		"can not increment projection offset for persistence key %s, offset %d is not the current offset",
		pk,
		c,
	)
}

// insertOffset inserts an entry for the saga pk with the offset set to 1.
// It returns false if an entry already exists.
func (OffsetStore) insertOffset(
	ctx context.Context,
	tx *sql.Tx,
	pk string,
) (bool, error) {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO ax_projection_offset SET
			persistence_key = ?,
			next_offset = 1`,
		pk,
	)

	if sqlutil.IsDuplicateEntry(err) {
		return false, nil
	}

	return true, err
}

// updateOffset increments the entry for the saga pk.
// It returns false if c is not the currently stored offset.
func (OffsetStore) updateOffset(
	ctx context.Context,
	tx *sql.Tx,
	pk string,
	c uint64,
) (bool, error) {
	var current uint64

	err := tx.QueryRowContext(
		ctx,
		`SELECT
			next_offset
		FROM ax_projection_offset
		WHERE persistence_key = ?
		FOR UPDATE`,
		pk,
	).Scan(
		&current,
	)

	if err != sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if c != current {
		return false, nil
	}

	return true, sqlutil.ExecSingleRow(
		ctx,
		tx,
		`UPDATE ax_projection_offset SET
			next_offset = next_offset + 1
		WHERE persistence_key = ?`,
		pk,
	)
}
