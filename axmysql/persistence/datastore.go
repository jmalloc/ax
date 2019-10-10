package persistence

import (
	"context"
	"database/sql"

	"github.com/jmalloc/ax/persistence"
)

// DataStore is a MySQL-backed implementation of Ax's persistence.DataStore
// interface.
type DataStore struct {
	DB *sql.DB
}

// BeginTx starts a new transaction.
func (ds *DataStore) BeginTx(ctx context.Context) (persistence.Tx, persistence.Committer, error) {
	tx, err := ds.DB.BeginTx(ctx, txOptions)
	if err != nil {
		return nil, nil, err
	}

	return &Tx{ds, tx}, tx, nil
}

// txOptions is the set of options used when starting a new SQL transaction.
var txOptions = &sql.TxOptions{
	Isolation: sql.LevelReadCommitted,
}

// ExtractDB returns the SQL database within ds.
// It panics if ds is not a *DataStore.
func ExtractDB(ds persistence.DataStore) *sql.DB {
	return ds.(*DataStore).DB
}
