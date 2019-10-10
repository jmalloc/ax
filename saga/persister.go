package saga

import (
	"context"

	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/persistence"
)

// Persister is an interface for loading saga instances, and persisting the
// changes that occur to them.
type Persister interface {
	// BeginUnitOfWork starts a new unit-of-work that modifies a saga instance.
	//
	// If the saga instance does not exist, it returns a UnitOfWork with an
	// instance at revision zero.
	BeginUnitOfWork(
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
	//
	// It returns true if any changes have occurred.
	// On success, the Instance().Revision is updated to match the new revision.
	Save(ctx context.Context) (bool, error)

	// SaveAndComplete persists changes to a completed instance.
	//
	// The precise behavior is implementation defined. Typically meta-data about
	// the instance is discarded. The implementation may completely remove any
	// record of the instance.
	//
	// On success, the Instance().Revision is updated to match the revision
	// produced by the save.
	SaveAndComplete(ctx context.Context) error

	// Close is called when the unit-of-work has ended, regardless of whether
	// Save() has been called.
	Close()
}
