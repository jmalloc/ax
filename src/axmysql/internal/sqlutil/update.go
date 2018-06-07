package sqlutil

import (
	"context"
	"database/sql"
	"fmt"
)

// UpdateSingleRow performs an update and verifies that only a single row was
// updated.
//
// q is the SQL update statement to execute. v is the variable arguments to pass
// along to the execute statement.
func UpdateSingleRow(
	ctx context.Context,
	tx *sql.Tx,
	q string,
	v ...interface{},
) error {
	res, err := tx.ExecContext(ctx, q, v...)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return fmt.Errorf("update to single row affected %d rows", n)
	}

	return nil
}
