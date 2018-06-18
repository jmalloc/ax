package persistence

import (
	"context"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// DataStore is a Bolt-backed implementation of Ax's persistence.DataStore
// interface.
type DataStore struct {
	DB *bolt.DB
}

// BeginTx starts a new transaction.
func (ds *DataStore) BeginTx(ctx context.Context) (persistence.Tx, persistence.Committer, error) {
	tx, err := ds.DB.Begin(true)
	if err != nil {
		return nil, nil, err
	}

	return &Tx{ds, tx}, tx, nil
}

// ExtractDB returns the Bolt database within ds.
// It panics if ds is not a *DataStore.
func ExtractDB(ds persistence.DataStore) *bolt.DB {
	return ds.(*DataStore).DB
}
