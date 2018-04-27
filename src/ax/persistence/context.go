package persistence

// contextKey is a type used for the keys of context values. A specific type is
// used to prevent collisions with context keys from other packages.
type contextKey string

const (
	dataStoreKey   contextKey = "ds"
	transactionKey contextKey = "tx"
)
