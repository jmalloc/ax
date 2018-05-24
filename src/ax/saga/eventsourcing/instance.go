package eventsourcing

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/ax/saga"
)

// Data is a specialization of Data for sagas that use eventsourcing.
type Data interface {
	saga.Data

	// ApplyEvent updates the data to reflect the fact that an event has
	// occurred.
	//
	// It may panic if env.Message does not implement ax.Event.
	ApplyEvent(env ax.Envelope)
}

// InstanceRepository is an interface for loading and saving eventsourced saga
// instances.
type InstanceRepository interface {
	// LoadSagaInstance fetches a saga instance by its ID.
	//
	// If a saga instance is found; ok is true, otherwise it is false. A
	// non-nil error indicates a problem with the store itself.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	LoadSagaInstance(
		ctx context.Context,
		tx persistence.Tx,
		id saga.InstanceID,
		d Data,
	) (saga.Instance, error)

	// SaveSagaInstance persists a saga instance.
	//
	// It returns an error if the saga instance has been modified since it was
	// loaded, or if there is a problem communicating with the store itself.
	//
	// It panics if envs contains any messages that do not implement ax.Event.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	SaveSagaInstance(
		ctx context.Context,
		tx persistence.Tx,
		i saga.Instance,
		envs []ax.Envelope,
	) error
}
