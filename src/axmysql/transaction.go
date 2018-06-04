package axmysql

import (
	"context"
	"database/sql"

	"github.com/jmalloc/ax/src/ax/persistence"
)

// Tx is an implementation of persistence.Tx that wraps an SQL transaction.
type Tx struct {
	ds    *DataStore
	sqlTx *sql.Tx
}

// DataStore returns the DataStore that the transaction operates on.
func (tx *Tx) DataStore() persistence.DataStore {
	return tx.ds
}

// GetTx returns the SQL transaction contained in ctx.
// It panics if ctx does not contain an SQL transaction.
func GetTx(ctx context.Context) *sql.Tx {
	tx, _ := persistence.GetTx(ctx)
	return sqlTx(tx)
}

// sqlTx returns the standard SQL transaction wrapped by tx.
// It panics if tx is not *axmysql.Tx
func sqlTx(tx persistence.Tx) *sql.Tx {
	return tx.(*Tx).sqlTx
}
