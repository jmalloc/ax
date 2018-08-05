package persistence

import (
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/protavo/src/protavo/driver"
)

// Tx is a Protavo-backed implementation of Ax's persistence.Tx interface.
type Tx struct {
	ds        *DataStore
	protavoTx driver.WriteTx
}

// DataStore returns the DataStore that the transaction operates on.
func (tx *Tx) DataStore() persistence.DataStore {
	return tx.ds
}

// ExtractTx returns the Protavo transaction within tx.
// It panics if tx is not a *Tx.
func ExtractTx(tx persistence.Tx) driver.WriteTx {
	return tx.(*Tx).protavoTx
}

type committer struct {
	protavoTx driver.WriteTx
}

func (c *committer) Commit() error {
	return c.protavoTx.Commit()
}

func (c *committer) Rollback() error {
	return c.protavoTx.Close()
}
