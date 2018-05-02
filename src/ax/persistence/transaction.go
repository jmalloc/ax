package persistence

import (
	"context"
	"errors"
)

// Tx represents an atomic unit-of-work performed on a DataStore.
type Tx interface {
	// DataStore returns the DataStore that the transaction operats on.
	DataStore() DataStore
}

// Committer is an interface used to commit and rollback persistence
// transactions.
type Committer interface {
	// Commit applies the changes to the data store.
	Commit() error

	// Rollback discards the changes without applying them to the data store.
	Rollback() error
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
func GetTx(ctx context.Context) (tx Tx, ok bool) {
	v := ctx.Value(transactionKey)

	if v != nil {
		tx = v.(Tx)
		ok = true
	}

	return
}

// GetOrBeginTx returns the transaction stored in ctx, or starts a new one if
// ctx does not contain a transaction.
//
// If a new transaction is started, the caller is said to "own" the transaction,
// that is, the caller is responsible for committing the transaction.
//
// If ctx already contains a transaction, the caller is said to "participate" in
// the transaction, but is not responsible for committing. In this case, the
// returned Committer is configured such that Commit() and Rollback() are
// no-ops that always return nil.
func GetOrBeginTx(ctx context.Context) (Tx, Committer, error) {
	if tx, ok := GetTx(ctx); ok {
		return tx, noOpCommitter{}, nil
	}

	return BeginTx(ctx)
}

// BeginTx starts a new transaction using the data store in ctx.
func BeginTx(ctx context.Context) (Tx, Committer, error) {
	if ds, ok := GetDataStore(ctx); ok {
		return ds.BeginTx(ctx)
	}

	return nil, nil, errors.New("can not begin transaction, no data store is available in ctx")
}

// noOpCommitter is an implementation of Committer has no-op commit and rollback
// operations.
type noOpCommitter struct{}

func (noOpCommitter) Commit() error   { return nil }
func (noOpCommitter) Rollback() error { return nil }
