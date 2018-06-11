package sqlutil

import (
	"context"
	"database/sql"
	"fmt"
)

// ExecSingleRow executes a query without returning any rows and verifies that
// exactly one row was affected.
//
// For an UPDATE query, this means that the query must actually alter the row,
// not simply match one row.
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

// ExecInsertOrUpdate executes a MySQL "INSERT ... ON DUPLICATE KEY UPDATE"
// query, and returns true if it results in an insert.
func ExecInsertOrUpdate(
	ctx context.Context,
	tx *sql.Tx,
	query string,
	args ...interface{},
) (bool, error) {
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return false, err
	}

	// https://dev.mysql.com/doc/refman/8.0/en/insert-on-duplicate.html
	//
	// From the MySQL docs:
	//
	// With ON DUPLICATE KEY UPDATE, the affected-rows value per row is 1 if the
	// row is inserted as a new row, 2 if an existing row is updated, and 0 if an
	// existing row is set to its current values.
	n, err := res.RowsAffected()

	return n == 1, err
}
