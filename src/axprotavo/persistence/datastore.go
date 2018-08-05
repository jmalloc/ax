package persistence

import (
	"context"

	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/protavo/src/protavo"
)

// DataStore is a Protavo-backed implementation of Ax's persistence.DataStore
// interface.
type DataStore struct {
	DB *protavo.DB
}

// BeginTx starts a new transaction.
func (ds *DataStore) BeginTx(ctx context.Context) (persistence.Tx, persistence.Committer, error) {
	tx, err := ds.DB.BeginWrite(ctx)
	if err != nil {
		return nil, nil, err
	}

	return &Tx{ds, tx}, &committer{tx}, nil
}

// ExtractDB returns the Protavo database within ds.
// It panics if ds is not a *DataStore.
func ExtractDB(ds persistence.DataStore) *protavo.DB {
	return ds.(*DataStore).DB
}
