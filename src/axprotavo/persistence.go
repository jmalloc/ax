package axprotavo

import (
	"context"

	"github.com/jmalloc/ax/src/ax/persistence"
	protavopersistence "github.com/jmalloc/ax/src/axprotavo/persistence"
	"github.com/jmalloc/protavo/src/protavo"
	"github.com/jmalloc/protavo/src/protavo/driver"
)

// NewDataStore returns a new data store that is backed by a Protavo database.
func NewDataStore(db *protavo.DB) persistence.DataStore {
	return &protavopersistence.DataStore{DB: db}
}

// GetDB returns the SQL database contained in ctx.
//
// It panics if ctx does not contain a Protavo-specific SQL database.
func GetDB(ctx context.Context) *protavo.DB {
	ds, _ := persistence.GetDataStore(ctx)
	return protavopersistence.ExtractDB(ds)
}

// GetTx returns the SQL transaction contained in ctx.
//
// It panics if ctx does not contain a Protavo-specific SQL transaction.
func GetTx(ctx context.Context) driver.WriteTx {
	tx, _ := persistence.GetTx(ctx)
	return protavopersistence.ExtractTx(tx)
}
