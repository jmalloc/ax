package persistence

import "context"

// DataStore is an interface for accessing a transactional data store.
type DataStore interface {
	// BeginTx starts a new transaction.
	BeginTx(ctx context.Context) (Tx, Committer, error)
}

// WithDataStore returns a new context derived from parent that contains
// a DataStore.
//
// The data store can be retrieved from the context with GetDataStore().
func WithDataStore(parent context.Context, ds DataStore) context.Context {
	return context.WithValue(parent, dataStoreKey, ds)
}

// GetDataStore returns the DataStore in ctx.
//
// If ctx does not contain a data store then ok is false.
func GetDataStore(ctx context.Context) (ds DataStore, ok bool) {
	v := ctx.Value(dataStoreKey)

	if v != nil {
		ds = v.(DataStore)
		ok = true
	}

	return
}
