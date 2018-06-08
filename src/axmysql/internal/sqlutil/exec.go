package sqlutil

import (
	"context"
	"database/sql"
	"fmt"
)

// ExecSingleRow executes a query without returning any rows and verifies that
// only a single row was affected.
//
// The args are for any placeholder parameters in the query. It returns an error
// if more than one row was affected.
func ExecSingleRow(
	ctx context.Context,
	tx *sql.Tx,
	query string,
	args ...interface{},
) error {
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return fmt.Errorf("execution of query on single row actually affected %d rows", n)
	}

	return nil
}
