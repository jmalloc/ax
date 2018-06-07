package endpoint

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// SinkSender is an implementation of ax.Sender that passes messages to a
// message sink.
type SinkSender struct {
	Sink       MessageSink
	Validators []Validator
}

// ExecuteCommand sends a command message.
//
// If ctx contains a message envelope, m is sent as a child of the message in
// that envelope.
func (s SinkSender) ExecuteCommand(ctx context.Context, m ax.Command) (ax.Envelope, error) {
	return s.send(ctx, OpSendUnicast, m)
}

// PublishEvent sends an event message.
//
// If ctx contains a message envelope, m is sent as a child of the message in
// that envelope.
func (s SinkSender) PublishEvent(ctx context.Context, m ax.Event) (ax.Envelope, error) {
	return s.send(ctx, OpSendMulticast, m)
}

// send wraps m in an envelope and passes that envelope to s.Sink.
// The new envelope is configured as a child of the envelope in ctx, if any.
func (s SinkSender) send(ctx context.Context, op Operation, m ax.Message) (ax.Envelope, error) {
	env, ok := GetEnvelope(ctx)

	if ok {
		env = env.NewChild(m)
	} else {
		env = ax.NewEnvelope(m)
	}

	if len(s.Validators) == 0 {
		s.Validators = DefaultValidators
	}

	for _, v := range s.Validators {
		if err := v.Validate(ctx, m); err != nil {
			return ax.Envelope{}, err
		}
	}

	return env, s.Sink.Accept(
		ctx,
		OutboundEnvelope{
			Envelope:  env,
			Operation: op,
		},
	)
}
