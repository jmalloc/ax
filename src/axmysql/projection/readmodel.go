package projection

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/persistence"
	mysqlpersistence "github.com/jmalloc/ax/src/axmysql/persistence"
	"github.com/jmalloc/ax/src/internal/typeswitch"
)

// ReadModel is an interface for application defined read-model projectors.
//
// Read-model projectors are a specialization of projectors that are designed to
// produce an application read-model from a series of events, and persist that
// read-model in a MySQL database.
//
// For each event type to be applied to the read-model, the projector must
// implement an "apply" method that adheres to one of the following signatures:
//
//     func (ctx context.Context, tx *sql.Tx, ev *<T>)
//     func (ctx context.Context, tx *sql.Tx, ev *<T>, env ax.Envelope)
//
// Where T is a struct type that implements ax.Event.
//
// Applier methods are responsible for mutating the read-model state. The
// appropriate applier is called for each message encountered in the stream. Any
// messages in the stream that do not have an associated applier method are
// ignored.
//
// The names of handler methods are meaningful. Each handler method's name must
// begin with "When". By convention these prefixes are followed by the message
// name, such as:
//
//     func (*BankAccount) WhenAccountCredited(*messages.AccountCredited)
type ReadModel interface {
	// PersistenceKey returns a unique name for the read-model.
	//
	// The persistence key is used to relate persisted data with the read-model
	// implementation that owns it. Persistence keys should not be changed once
	// the read-model's projector has been started.
	PersistenceKey() string
}

// ReadModelProjector is a projector that applies events to a ReadModel.
type ReadModelProjector struct {
	ReadModel  ReadModel
	EventTypes ax.MessageTypeSet
	Apply      typeswitch.Switch
}

// NewReadModelProjector returns a new projector that applies events to a
// read-model.
func NewReadModelProjector(rm ReadModel) *ReadModelProjector {
	p := &ReadModelProjector{
		ReadModel: rm,
	}

	sw, _, err := typeswitch.New(
		[]reflect.Type{
			reflect.TypeOf(rm),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
			reflect.TypeOf((*ax.Envelope)(nil)).Elem(),
			reflect.TypeOf((*context.Context)(nil)).Elem(),
			reflect.TypeOf((*sql.Tx)(nil)),
		},
		[]reflect.Type{
			reflect.TypeOf((*error)(nil)).Elem(),
		},
		readModelApplySignature,
		readModelApplyWithEnvelopeSignature,
	)
	if err != nil {
		panic(err)
	}

	p.Apply = sw
	p.EventTypes = ax.TypesByGoType(sw.Types()...)

	return p
}

// PersistenceKey returns a unique name for the projector.
//
// The persistence key is used to relate persisted data with the projector
// implementation that owns it. Persistence keys should not be changed once
// a projection has been started.
func (p ReadModelProjector) PersistenceKey() string {
	return p.ReadModel.PersistenceKey()
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

	out := p.Apply.Dispatch(
		p.ReadModel,
		env.Message.(ax.Event),
		env,
		ctx,
		tx,
	)

	if err := out[0]; err != nil {
		return err.(error)
	}

	return nil
}

var (
	readModelApplySignature = &typeswitch.Signature{
		In: []reflect.Type{
			reflect.TypeOf((*ReadModel)(nil)).Elem(),
			reflect.TypeOf((*context.Context)(nil)).Elem(),
			reflect.TypeOf((*sql.Tx)(nil)),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
		},
		Out: []reflect.Type{
			reflect.TypeOf((*error)(nil)).Elem(),
		},
	}

	readModelApplyWithEnvelopeSignature = &typeswitch.Signature{
		In: []reflect.Type{
			reflect.TypeOf((*ReadModel)(nil)).Elem(),
			reflect.TypeOf((*context.Context)(nil)).Elem(),
			reflect.TypeOf((*sql.Tx)(nil)),
			reflect.TypeOf((*ax.Event)(nil)).Elem(),
			reflect.TypeOf((*ax.Envelope)(nil)).Elem(),
		},
		Out: []reflect.Type{
			reflect.TypeOf((*error)(nil)).Elem(),
		},
	}
)
