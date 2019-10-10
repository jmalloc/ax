package axdogma

import (
	"context"
	"errors"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/projection"
)

// ProjectionAdaptor adapts a dogma.ProjectionMessageHandler to Ax's
// projection.Projector interface.
type ProjectionAdaptor struct {
	Key        string
	EventTypes ax.MessageTypeSet
	Handler    dogma.ProjectionMessageHandler
}

var _ projection.Projector = &ProjectionAdaptor{}

// PersistenceKey returns a unique name for the projector.
//
// The persistence key is used to relate persisted data with the projector
// implementation that owns it. Persistence keys should not be changed once
// a projection has been started.
func (a *ProjectionAdaptor) PersistenceKey() string {
	return a.Key
}

// MessageTypes returns the set of messages that the projector intends
// to handle.
//
// The return value should be constant as it may be cached.
func (a *ProjectionAdaptor) MessageTypes() ax.MessageTypeSet {
	return a.EventTypes
}

// ApplyMessage invokes application-defined logic that updates the
// application state to reflect the occurrence of a message.
//
// It may panic if env.Message is not one of the types described by
// MessageTypes().
func (a *ProjectionAdaptor) ApplyMessage(ctx context.Context, mctx ax.MessageContext) error {
	res, err := mctx.Envelope.MessageID.MarshalText()
	if err != nil {
		return err
	}

	cur, err := a.Handler.ResourceVersion(ctx, res)
	if err != nil {
		return err
	}

	// This message has already been handled.
	if len(cur) != 0 {
		return nil
	}

	ok, err := a.Handler.HandleEvent(
		ctx,
		res, cur, []byte{1},
		&projectionScope{mctx},
		mctx.Envelope.Message,
	)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	return errors.New("optimsitic concurrency conflict occurred while updating projection")
}

type projectionScope struct {
	mctx ax.MessageContext
}

func (s *projectionScope) RecordedAt() time.Time {
	return s.mctx.Envelope.SendAt
}

func (s *projectionScope) Log(f string, v ...interface{}) {
	s.mctx.Log(f, v...)
}
