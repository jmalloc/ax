package logging

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/twelf/src/twelf"
)

// InboundStage is an inbound pipeline stage that logs information about
// inbound messages.
type InboundStage struct {
	Logger twelf.Logger
	Next   bus.InboundPipeline
}

// Initialize is called after the transport is initialized. It can be used
// to inspect or configure the transport as per the needs of the pipeline.
func (l *InboundStage) Initialize(ctx context.Context, t bus.Transport) error {
	return l.Next.Initialize(ctx, t)
}

// Accept forwards an inbound message through the pipeline until
// it is handled by some application-defined message handler(s).
func (l *InboundStage) Accept(ctx context.Context, s bus.MessageSink, env bus.InboundEnvelope) error {
	t := time.Now()
	if err := l.Next.Accept(ctx, s, env); err != nil {
		return err
	}

	d := time.Since(t)
	mt := env.Type()

	l.Logger.Log(
		"handled %s %s: %s (%s) [%s, %s]",
		formatMessageType(mt),
		env.MessageID,
		env.Message.Description(),
		d,
		env.CausationID,
		env.CorrelationID,
	)

	return nil
}
