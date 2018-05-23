package axmysql

import (
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

// sqlTx returns the standard SQL transaction wrapped by tx.
// It panics if tx is not *axmysql.Tx
func sqlTx(tx persistence.Tx) *sql.Tx {
	return tx.(*Tx).sqlTx
}
