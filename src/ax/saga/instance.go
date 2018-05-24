package saga

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax/ident"
	"github.com/jmalloc/ax/src/ax/persistence"
)

// InstanceID uniquely identifies a saga instance.
type InstanceID struct {
	ident.ID
}

// Instance is an instance of a saga.
//
// It encapsulates the application-defined saga data and its meta-data.
type Instance struct {
	// InstanceID is a globally unique identifier for the saga instance.
	InstanceID InstanceID

	// Data is the application-defined data associated with this instance.
	Data Data

	// Revision is version of the instance that the data represents.
	Revision Revision
}

// Revision is the version of a saga instance.
type Revision uint64

// Data is an interface for application-defined data associated with a saga
// instance.
type Data interface {
	proto.Message

	// SagaDescription returns a human-readable description of the saga
	// instance.
	//
	// Assume that the description will be used inside log messages or displayed
	// in audit logs.
	//
	// Follow the same conventions as for error messages:
	// https://github.com/golang/go/wiki/CodeReviewComments#error-strings
	SagaDescription() string
}

// InstanceRepository is an interface for loading and saving saga instances.
type InstanceRepository interface {
	// LoadSagaInstance fetches a saga instance by its ID.
	//
	// If a saga instance is found; ok is true, otherwise it is false. A
	// non-nil error indicates a problem with the store itself.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	LoadSagaInstance(ctx context.Context, tx persistence.Tx, id InstanceID) (Instance, error)

	// SaveSagaInstance persists a saga instance.
	//
	// It returns an error if the saga instance has been modified since it was
	// loaded, or if there is a problem communicating with the store itself.
	//
	// It panics if the repository is not able to enlist in tx because it uses a
	// different underlying storage system.
	SaveSagaInstance(ctx context.Context, tx persistence.Tx, i Instance) error
}
