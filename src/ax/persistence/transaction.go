package persistence

import (
	"context"
	"errors"
)

// Tx represents an atomic unit-of-work performed on a DataStore.
type Tx interface {
	// Commit applies the changes to the data store.
	Commit() error

	// Rollback discards the changes without applying them to the data store.
	Rollback() error

	// UnderlyingTx returns a data-store-specific value that represents the
	// transaction. For example, for an SQL-based data store this may be the
	// *sql.Tx.
	UnderlyingTx() interface{}
}

// WithTx returns a new context derived from parent that contains a transaction.
//
// The transaction can be retrieved from the context with GetTx().
func WithTx(parent context.Context, tx Tx) context.Context {
	return context.WithValue(parent, transactionKey, tx)
}

// GetTx returns the transaction stored in ctx.
//
// If ctx does not contain a transaction then ok is false.
//
// Transactions are made available via the context so that application-defined
// message handlers can optionally perform some additional storage operations
// within the same transaction as infrastructure features such as the outbox
// system.
//
// Care should be taken not to commit or rollback the transaction within
// a message handler.
func GetTx(ctx context.Context) (tx Tx, ok bool) {
	v := ctx.Value(transactionKey)

	if v != nil {
		tx = v.(Tx)
		ok = true
	}

	return
}

// GetOrBeginTx returns the transaction stored in ctx, or starts a new one
// ctx does not contain a transaction.
//
// If a new transaction is started, the caller is said to "own" the transaction,
// that is, the caller is responsible for committing the transaction.
//
// If ctx contains a transaction, the caller is said to "participate" in the
// transaction, but is not responsible for committing. In this case, the
// returned transaction is configured such that Commit() and Rollback() are
// no-ops that always return nil.
func GetOrBeginTx(ctx context.Context) (Tx, error) {
	if tx, ok := GetTx(ctx); ok {
		return noOpTx{tx}, nil
	}

	return BeginTx(ctx)
}

// BeginTx starts a new transactions using the data store in ctx.
func BeginTx(ctx context.Context) (Tx, error) {
	if ds, ok := GetDataStore(ctx); ok {
		return ds.BeginTx(ctx)
	}

	return nil, errors.New("can not begin transaction, no data store is available in ctx")
}

// noOpTx is a transaction wrapper that has no-op commit and rollback operations.
type noOpTx struct {
	tx Tx
}

func (tx noOpTx) Commit() error             { return nil }
func (tx noOpTx) Rollback() error           { return nil }
func (tx noOpTx) UnderlyingTx() interface{} { return tx.tx.UnderlyingTx() }
