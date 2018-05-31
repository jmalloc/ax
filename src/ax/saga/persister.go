package saga

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// Persister is an interface for loading saga instances, and persisting the
// changes that occur to them.
type Persister interface {
	// BeginCreate starts a new unit-of-work that persists a new saga instance.
	BeginCreate(
		ctx context.Context,
		sg Saga,
		tx persistence.Tx,
		s ax.Sender,
		i Instance,
	) (UnitOfWork, error)

	// BeginUpdate starts a new unit-of-work that updates an existing saga
	// instance.
	BeginUpdate(
		ctx context.Context,
		sg Saga,
		tx persistence.Tx,
		s ax.Sender,
		id InstanceID,
	) (UnitOfWork, error)
}

// UnitOfWork encapsulates the logic for persisting changes to an instance.
type UnitOfWork interface {
	// Sender returns the ax.Sender that the saga must use to send messages.
	// This allows the persister to capture new messages if necessary.
	Sender() ax.Sender

	// Instance returns the saga instance that the unit-of-work applies to.
	Instance() Instance

	// Save persists changes to the instance.
	// It returns true if any changes have occurred.
	Save(ctx context.Context) (bool, error)
}
