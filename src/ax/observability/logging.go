package observability

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/twelf/src/twelf"
)

// LoggingObserver is an observer that logs about messages.
type LoggingObserver struct {
	Logger twelf.Logger
}

// BeforeInbound logs information about an inbound message.
func (o *LoggingObserver) BeforeInbound(ctx context.Context, env endpoint.InboundEnvelope) {
	o.logMessage(ctx, "▼", env.Envelope)
}

// AfterInbound logs information about errors that occur processing an inbound message.
func (o *LoggingObserver) AfterInbound(ctx context.Context, env endpoint.InboundEnvelope, err error) {
	if err != nil {
		o.logError(ctx, "▽", env.Envelope, err)
	}
}

// BeforeOutbound logs information about an outbound message.
func (o *LoggingObserver) BeforeOutbound(ctx context.Context, env endpoint.OutboundEnvelope) {
	o.logMessage(ctx, "▲", env.Envelope)
}

// AfterOutbound logs information about an outbound message.
func (o *LoggingObserver) AfterOutbound(ctx context.Context, env endpoint.OutboundEnvelope, err error) {
	if err != nil {
		o.logError(ctx, "△", env.Envelope, err)
	}
}

// log writes a message to o.Logger. If o.Logger is nil, it uses twelf.DefaultLogger.
func (o *LoggingObserver) log(s string) {
	l := o.Logger

	if l == nil {
		l = twelf.DefaultLogger
	}

	l.LogString(s)
}

func (o *LoggingObserver) logMessage(ctx context.Context, icon string, env ax.Envelope) {
	s := fmt.Sprintf(
		"%s   %s  %s",
		icon,
		env.Message.MessageDescription(),
		formatEnvelope(env),
	)

	if in, ok := endpoint.GetEnvelope(ctx); ok {
		s += " " + formatAttempt(in)
	}

	o.log(s)
}

func (o *LoggingObserver) logError(ctx context.Context, icon string, env ax.Envelope, err error) {
	s := fmt.Sprintf(
		"%s ✘ %s ∎ %s  %s",
		icon,
		env.Message.MessageDescription(),
		err,
		formatEnvelope(env),
	)

	if in, ok := endpoint.GetEnvelope(ctx); ok {
		s += " " + formatAttempt(in)
	}

	o.log(s)
}

func formatEnvelope(env ax.Envelope) string {
	return fmt.Sprintf(
		"[%s msg:%s cause:%s corr:%s]",
		env.Type(),
		env.MessageID,
		env.CausationID,
		env.CorrelationID,
	)
}

func formatAttempt(env endpoint.InboundEnvelope) string {
	return fmt.Sprintf(
		"[attempt:%s #%d]",
		env.AttemptID,
		env.AttemptCount,
	)
}
