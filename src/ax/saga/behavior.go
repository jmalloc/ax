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
func (IgnoreNotFound) HandleNotFound(context.Context, ax.Sender, ax.MessageContext) error {
	return nil
}

// ErrorIfNotFound is an embeddable struct that implements a
// Saga.HandleNotFound() method that always returns an error.
type ErrorIfNotFound struct{}

// HandleNotFound always returns an error.
func (ErrorIfNotFound) HandleNotFound(context.Context, ax.Sender, ax.MessageContext) error {
	return errors.New("could not find a saga instance to handle message")
}

// CompletableByData is an embeddable struct that implements a
// Saga.IsInstanceComplete() method that forwards the completion check on to a
// CompletableData value.
type CompletableByData struct{}

// IsInstanceComplete returns true if i.Data implements CompletableData and
// i.Data.IsInstanceComplete() returns true.
func (CompletableByData) IsInstanceComplete(_ context.Context, i Instance) (bool, error) {
	if cd, ok := i.Data.(CompletableData); ok {
		return cd.IsInstanceComplete(), nil
	}

	return false, nil
}

// InstancesNeverComplete is an embeddable struct that implements a
// Saga.IsInstanceComplete() method that always returns false.
type InstancesNeverComplete struct{}

// IsInstanceComplete always returns false.
func (InstancesNeverComplete) IsInstanceComplete(context.Context, Instance) (bool, error) {
	return false, nil
}
