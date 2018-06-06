package projection

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmalloc/ax/src/ax/persistence"
	mysqlpersistence "github.com/jmalloc/ax/src/axmysql/persistence"
)

// OffsetStore is a MySQL-backed implementation of Ax's projection.OffsetStore
// interface.
type OffsetStore struct{}

// LoadOffset returns the offset at which a consumer should resume
// reading from the stream.
//
// pn is the projector name.
func (OffsetStore) LoadOffset(
	ctx context.Context,
	ds persistence.DataStore,
	pn string,
) (uint64, error) {
	db := mysqlpersistence.ExtractDB(ds)

	var offset uint64

	err := db.QueryRowContext(
		ctx,
		`SELECT
			next_offset
		FROM projection_offset
		WHERE projection = ?`,
		pn,
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

// SaveOffset stores the next offset at which a consumer should resume
// reading from the stream.
//
// pn is the projector name. c is the offset that is currently stored, as
// returned by LoadOffset(). If c is not the offset that is currently stored,
// a non-nil error is returned. o is the new offset to store.
func (OffsetStore) SaveOffset(
	ctx context.Context,
	ptx persistence.Tx,
	pn string,
	c, o uint64,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

	if c == 0 {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO projection_offset SET
				projection = ?,
				next_offset = ?`,
			pn,
			o,
		)

		return err
	}

	res, err := tx.ExecContext(
		ctx,
		`UPDATE projection_offset SET
			next_offset = ?
		WHERE projection = ?
		AND next_offset = ?`,
		o,
		pn,
		c,
	)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf(
			"can not store offset for %s projection, offset %d is not the currently stored offset",
			pn,
			c,
		)
	}

	return err
}
