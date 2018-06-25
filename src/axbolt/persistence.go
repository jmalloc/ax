package axbolt

import (
	"context"

	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/ax/src/ax/persistence"
	boltpersistence "github.com/jmalloc/ax/src/axbolt/persistence"
)

// NewDataStore returns a new data store that is backed by a Bolt database.
func NewDataStore(db *bolt.DB) persistence.DataStore {
	return &boltpersistence.DataStore{DB: db}
}

// GetDB returns the database contained in ctx.
//
// It panics if ctx does not contain a Bolt-specific database.
func GetDB(ctx context.Context) *bolt.DB {
	ds, _ := persistence.GetDataStore(ctx)
	return boltpersistence.ExtractDB(ds)
}

// GetTx returns the Bolt-specific transaction contained in ctx.
//
// It panics if ctx does not contain a Bolt-specific transaction.
func GetTx(ctx context.Context) *bolt.Tx {
	tx, _ := persistence.GetTx(ctx)
	return boltpersistence.ExtractTx(tx)
}
