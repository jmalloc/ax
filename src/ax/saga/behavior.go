package saga

import (
	"context"
	"errors"

	"github.com/jmalloc/ax/src/ax"
)

// IgnoreNotFound is an embeddable struct that implements a
// Saga.HandleNotFound() method that is a no-op.
type IgnoreNotFound struct{}

// HandleNotFound always returns nil.
func (IgnoreNotFound) HandleNotFound(context.Context, ax.Sender, ax.Envelope) error {
	return nil
}

// ErrorIfNotFound is an embeddable struct that implements a
// Saga.HandleNotFound() method that always returns an error.
type ErrorIfNotFound struct{}

// HandleNotFound always returns an error.
func (ErrorIfNotFound) HandleNotFound(_ context.Context, _ ax.Sender, _ ax.Envelope) error {
	return errors.New("could not find a saga instance to handle message")
}

// CompletableByData is an embeddable struct implements a
// CompletableSaga.IsComplete() method that forwards the completion check on to
// a CompletableData value.
type CompletableByData struct{}

// IsComplete returns true if i.Data implements CompletableData and
// i.Data.IsComplete() returns true.
func (CompletableByData) IsComplete(_ context.Context, i Instance) (bool, error) {
	if cd, ok := i.Data.(CompletableData); ok {
		return cd.IsComplete(), nil
	}

	return false, nil
}
