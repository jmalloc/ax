package axmysql

import (
	"context"
	"database/sql"

	mysqlpersistence "github.com/jmalloc/ax/axmysql/persistence"
	"github.com/jmalloc/ax/persistence"
)

// NewDataStore returns a new data store that is backed by a MySQL database.
func NewDataStore(db *sql.DB) persistence.DataStore {
	return &mysqlpersistence.DataStore{DB: db}
}

// GetDB returns the SQL database contained in ctx.
//
// It panics if ctx does not contain a MySQL-specific SQL database.
func GetDB(ctx context.Context) *sql.DB {
	ds, _ := persistence.GetDataStore(ctx)
	return mysqlpersistence.ExtractDB(ds)
}

// GetTx returns the SQL transaction contained in ctx.
//
// It panics if ctx does not contain a MySQL-specific SQL transaction.
func GetTx(ctx context.Context) *sql.Tx {
	tx, _ := persistence.GetTx(ctx)
	return mysqlpersistence.ExtractTx(tx)
}
