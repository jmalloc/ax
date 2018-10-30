package observability

import (
	"context"

	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/twelf/src/twelf"
)

// LoggingObserver is an observer that logs about messages.
type LoggingObserver struct {
	Logger twelf.Logger
}

// BeforeInbound logs information about an inbound message.
func (o *LoggingObserver) BeforeInbound(ctx context.Context, env endpoint.InboundEnvelope) {
	o.log(
		"recv: %s  [%s msg:%s cause:%s corr:%s del:%s#%d]",
		env.Message.MessageDescription(),
		env.Type(),
		env.MessageID,
		env.CausationID,
		env.CorrelationID,
		env.DeliveryID,
		env.DeliveryCount,
	)
}

// AfterInbound logs information about errors that occur processing an inbound message.
func (o *LoggingObserver) AfterInbound(ctx context.Context, env endpoint.InboundEnvelope, err error) {
	if err != nil {
		o.log(
			"recv error: %s  %s  [%s msg:%s cause:%s corr:%s del:%s#%d]",
			env.Message.MessageDescription(),
			err,
			env.Type(),
			env.MessageID,
			env.CausationID,
			env.CorrelationID,
			env.DeliveryID,
			env.DeliveryCount,
		)
	}
}

// BeforeOutbound logs information about an outbound message.
func (o *LoggingObserver) BeforeOutbound(ctx context.Context, env endpoint.OutboundEnvelope) {
	o.log(
		"send: %s  [%s msg:%s cause:%s corr:%s]",
		env.Message.MessageDescription(),
		env.Type(),
		env.MessageID,
		env.CausationID,
		env.CorrelationID,
	)
}

// AfterOutbound logs information about an outbound message.
func (o *LoggingObserver) AfterOutbound(ctx context.Context, env endpoint.OutboundEnvelope, err error) {
	if err != nil {
		o.log(
			"send error: %s  %s  [%s msg:%s cause:%s corr:%s]",
			env.Message.MessageDescription(),
			err,
			env.Type(),
			env.MessageID,
			env.CausationID,
			env.CorrelationID,
		)
	}
}

// log writes a message to o.Logger. If o.Logger is nil, it uses twelf.DefaultLogger.
func (o *LoggingObserver) log(f string, v ...interface{}) {
	l := o.Logger

	if l == nil {
		l = twelf.DefaultLogger
	}

	l.Log(f, v...)
}
