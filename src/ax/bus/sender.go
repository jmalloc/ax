package bus

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// SinkSender is an implementation of ax.Sender that passes messages to a
// message sink.
type SinkSender struct {
	Sink MessageSink
}

// ExecuteCommand sends a command message.
//
// Commands are routed to a single endpoint as per the routing rules of the
// outbound message pipeline.
func (s SinkSender) ExecuteCommand(ctx context.Context, m ax.Command) error {
	return s.send(ctx, OpSendUnicast, m)
}

// PublishEvent sends an event message.
//
// Events are routed to endpoints that subscribe to messages of that type.
func (s SinkSender) PublishEvent(ctx context.Context, m ax.Event) error {
	return s.send(ctx, OpSendMulticast, m)
}

// send wraps m in an envelope and passes that envelope to s.Sink.
// The new envelope is configured as a child of the envelope in ctx, if any.
func (s SinkSender) send(ctx context.Context, op OutboundOperation, m ax.Message) error {
	env, ok := GetEnvelope(ctx)

	if ok {
		env = env.NewChild(m)
	} else {
		env = ax.NewEnvelope(m)
	}

	return s.Sink.Accept(
		ctx,
		OutboundEnvelope{
			Envelope:  env,
			Operation: op,
		},
	)
}
