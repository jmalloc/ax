package observability

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/twelf/src/twelf"
)

// LoggingObserver is an observer that logs about messages.
type LoggingObserver struct {
	Logger twelf.Logger
}

// BeforeInbound logs information about an inbound message.
func (o *LoggingObserver) BeforeInbound(ctx context.Context, env ax.Envelope) (context.Context, error) {
	mt := env.Type()

	o.log(
		"recv: %s  [%s msg:%s cause:%s corr:%s]",
		env.Message.Description(),
		mt,
		env.MessageID,
		env.CausationID,
		env.CorrelationID,
	)

	return ctx, nil
}

// AfterInbound logs information about errors that occur processing an inbound message.
func (o *LoggingObserver) AfterInbound(ctx context.Context, env ax.Envelope, acceptErr error) error {
	if acceptErr == nil {
		return nil
	}

	o.log(
		"error: %s  [%s msg:%s cause:%s corr:%s]",
		env.Message.Description(),
		acceptErr,
		env.Type(),
		env.MessageID,
		env.CausationID,
		env.CorrelationID,
	)

	return nil
}

// AfterOutbound logs information about an outbound message.
func (o *LoggingObserver) AfterOutbound(ctx context.Context, env ax.Envelope, acceptErr error) error {
	if acceptErr != nil {
		return nil
	}

	o.log(
		"send: %s  [%s msg:%s cause:%s corr:%s]",
		env.Message.Description(),
		env.Type(),
		env.MessageID,
		env.CausationID,
		env.CorrelationID,
	)

	return nil
}

func (o *LoggingObserver) log(f string, v ...interface{}) {
	l := o.Logger

	if l == nil {
		l = twelf.DefaultLogger
	}

	l.Log(f, v...)
}
