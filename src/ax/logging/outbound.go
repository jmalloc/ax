package logging

import (
	"context"

	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/twelf/src/twelf"
)

// OutboundStage is an outbound pipeline stage that logs information about
// outbound messages.
type OutboundStage struct {
	Logger twelf.Logger
	Next   bus.OutboundPipeline
}

// Initialize is called after the transport is initialized. It can be used
// to inspect or configure the transport as per the needs of the pipeline.
func (l *OutboundStage) Initialize(ctx context.Context, t bus.Transport) error {
	return l.Next.Initialize(ctx, t)
}

// Accept processes the message encapsulated in env.
func (l *OutboundStage) Accept(ctx context.Context, env bus.OutboundEnvelope) error {
	if err := l.Next.Accept(ctx, env); err != nil {
		return err
	}

	mt := env.Type()

	l.Logger.Log(
		"sent %s %s: %s [%s, %s]",
		formatMessageType(mt),
		env.MessageID,
		env.Message.Description(),
		env.CausationID,
		env.CorrelationID,
	)

	return nil
}
