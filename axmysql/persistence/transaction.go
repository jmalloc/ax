package persistence

import (
	"database/sql"

	"github.com/jmalloc/ax/persistence"
)

// Tx is a MySQL-backed implementation of Ax's persistence.Tx interface.
type Tx struct {
	ds    *DataStore
	sqlTx *sql.Tx
}

// DataStore returns the DataStore that the transaction operates on.
func (tx *Tx) DataStore() persistence.DataStore {
	return tx.ds
}

// ExtractTx returns the SQL transaction within tx.
// It panics if tx is not a *Tx.
func ExtractTx(tx persistence.Tx) *sql.Tx {
	return tx.(*Tx).sqlTx
}
