package projection

import (
	"context"
	"database/sql"

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
		FROM ax_projection_offset
		WHERE persistence_key = ?`,
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
			`INSERT INTO ax_projection_offset SET
				persistence_key = ?,
				next_offset = ?`,
			pn,
			o,
		)

		return err
	}

	return sqlutil.UpdateSingleRow(
		ctx,
		tx,
		`UPDATE ax_projection_offset SET
			next_offset = ?
		WHERE persistence_key = ?
		AND next_offset = ?`,
		o,
		pn,
		c,
	)
}
