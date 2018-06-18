package persistence

import (
	bolt "github.com/coreos/bbolt"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// Tx is a Bolt-backed implementation of Ax's persistence.Tx interface.
type Tx struct {
	ds     *DataStore
	boltTx *bolt.Tx
}

// DataStore returns the DataStore that the transaction operates on.
func (tx *Tx) DataStore() persistence.DataStore {
	return tx.ds
}

// ExtractTx returns the Bolt transaction within tx.
// It panics if tx is not a *Tx.
func ExtractTx(tx persistence.Tx) *bolt.Tx {
	return tx.(*Tx).boltTx
}
