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

// SaveOffset stores the next offset at which a consumer should resume
// reading from the stream.
//
// pk is the projector's persitence key. c is the offset that is currently
// stored, as returned by LoadOffset(). If c is not the offset that is
// currently stored, a non-nil error is returned. o is the new offset to store.
func (OffsetStore) SaveOffset(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	c, o uint64,
) error {
	tx := mysqlpersistence.ExtractTx(ptx)

	if c == 0 {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO ax_projection_offset SET
				persistence_key = ?,
				next_offset = ?`,
			pk,
			o,
		)

		return err
	}

	var no uint64
	err := tx.QueryRowContext(
		ctx,
		`SELECT
			next_offset
		FROM ax_projection_offset
		WHERE persistence_key = ?
		FOR UPDATE`,
		pk,
	).Scan(
		&no,
	)
	if err != nil {
		return err
	}

	if o != no {
		return fmt.Errorf(
			"can not update projection offset for persistence key %s, offset %d is not the next offset",
			pk,
			o,
		)
	}

	return sqlutil.ExecSingleRow(
		ctx,
		tx,
		`UPDATE ax_projection_offset SET
			next_offset = ?
		WHERE persistence_key = ?
		AND next_offset = ?`,
		o,
		pk,
		c,
	)
}
