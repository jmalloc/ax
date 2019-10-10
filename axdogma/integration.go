package axdogma

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/routing"
)

// IntegrationAdaptor adapts a dogma.IntegrationMessageHandler to Ax's
// routing.MessageHandler interface.
type IntegrationAdaptor struct {
	CommandTypes ax.MessageTypeSet
	Handler      dogma.IntegrationMessageHandler
}

var _ routing.MessageHandler = &IntegrationAdaptor{}

// MessageTypes returns the set of messages that the handler intends
// to handle.
//
// The return value should be constant as it may be cached by various
// independent stages in the message pipeline.
func (a *IntegrationAdaptor) MessageTypes() ax.MessageTypeSet {
	return a.CommandTypes
}

// HandleMessage invokes application-defined logic that handles a
// message.
//
// It may panic if env.Message is not one of the types described by
// MessageTypes().
func (a *IntegrationAdaptor) HandleMessage(
	ctx context.Context,
	s ax.Sender,
	mctx ax.MessageContext,
) error {
	sc := &integrationScope{
		mctx: mctx,
	}

	if err := a.Handler.HandleCommand(
		ctx,
		sc,
		mctx.Envelope.Message,
	); err != nil {
		return err
	}

	for _, m := range sc.events {
		if _, err := s.PublishEvent(ctx, m); err != nil {
			return err
		}
	}

	return nil
}

type integrationScope struct {
	mctx   ax.MessageContext
	events []ax.Event
}

func (s *integrationScope) RecordEvent(m dogma.Message) {
	s.events = append(s.events, m.(ax.Event))
}

func (s *integrationScope) Log(f string, v ...interface{}) {
	s.mctx.Log(f, v...)
}
