package projection

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	mysqlpersistence "github.com/jmalloc/ax/src/axmysql/persistence"
	"github.com/jmalloc/ax/src/internal/visitor"
)

// ReadModel is an interface for application defined read-model projectors.
//
// Read-model projectors are a specialization of projectors that are designed to
// produce an application read-model from a series of events, and persist that
// read-model in a MySQL database.
//
// For each event type to be applied to the read-model, the projector must
// implement an "apply" method that adheres to the following signature:
//
//     func (ctx context.Context, tx *sql.Tx, ev *<T>)
//
// Where T is a struct type that implements ax.Event.
//
// Applier methods are responsible for mutating the read-model state. The
// appropriate applier is called for each message encountered in the stream. Any
// messages in the stream that do not have an associated applier method are
// ignored.
//
// The names of applier methods are not meaningful to the projection system. By
// convention, event appliers are prefixed with the word "When", such as:
//
//     func (*BankAccount) WhenAccountCredited(*messages.AccountCredited)
type ReadModel interface {
	// ReadModelName returns a unique name for the read-model.
	//
	// The read-models's name is used to correlate persisted data with this
	// instance, so it should not be changed once data has been written.
	ReadModelName() string
}

// ReadModelProjector is a projector that applies events to a ReadModel.
type ReadModelProjector struct {
	ReadModel  ReadModel
	EventTypes ax.MessageTypeSet
	Apply      func(ReadModel, context.Context, *sql.Tx, ax.Event) error
}

// NewReadModelProjector returns a new projector that applies events to a
// read-model.
func NewReadModelProjector(rm ReadModel) *ReadModelProjector {
	p := &ReadModelProjector{
		ReadModel: rm,
	}

	eventTypes := visitor.MakeAcceptor(
		&p.Apply,
		reflect.TypeOf((*ax.Event)(nil)).Elem(),
		reflect.TypeOf(rm),
	)

	// TODO: make use of https://github.com/jmalloc/ax/issues/74
	for _, t := range eventTypes {
		p.EventTypes = p.EventTypes.Add(
			ax.TypeOf(
				reflect.Zero(t).Interface().(ax.Message),
			),
		)
	}

	return p
}

// ProjectorName returns a unique name for the projector.
//
// The projector's name is used to correlate persisted data with this
// instance, so it should not be changed.
func (p ReadModelProjector) ProjectorName() string {
	return p.ReadModel.ReadModelName()
}

// MessageTypes returns the set of messages that the projector intends
// to handle.
//
// The return value should be constant as it may be cached.
func (p ReadModelProjector) MessageTypes() ax.MessageTypeSet {
	return p.EventTypes
}

// ApplyMessage invokes application-defined logic that updates the
// application state to reflect the delivery of a message.
//
// It may panic if env.Message is not one of the types described by
// MessageTypes().
func (p ReadModelProjector) ApplyMessage(ctx context.Context, env ax.Envelope) error {
	ptx, _ := persistence.GetTx(ctx)
	tx := mysqlpersistence.ExtractTx(ptx)

	return p.Apply(
		p.ReadModel,
		ctx,
		tx,
		env.Message.(ax.Event),
	)
}
