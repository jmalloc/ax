package persistence

import (
	"database/sql"

	"github.com/jmalloc/ax/src/ax/persistence"
)

// ExtractDB returns the SQL database within ds.
// It panics if ds is not a *DataStore.
func ExtractDB(ds persistence.DataStore) *sql.DB {
	return ds.(*DataStore).DB
}

// ExtractTx returns the SQL transaction within tx.
// It panics if tx is not a *Tx.
func ExtractTx(tx persistence.Tx) *sql.Tx {
	return tx.(*Tx).sqlTx
}
