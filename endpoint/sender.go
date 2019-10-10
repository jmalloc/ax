package endpoint

import (
	"context"

	"github.com/jmalloc/ax"
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
func (s SinkSender) ExecuteCommand(
	ctx context.Context,
	m ax.Command,
	opts ...ax.ExecuteOption,
) (ax.Envelope, error) {
	env, err := s.newEnvelope(ctx, m)
	if err != nil {
		return ax.Envelope{}, err
	}

	for _, o := range opts {
		if err := o.ApplyExecuteOption(&env); err != nil {
			return ax.Envelope{}, err
		}
	}

	return env, s.Sink.Accept(
		ctx,
		OutboundEnvelope{
			Envelope:  env,
			Operation: OpSendUnicast,
		},
	)
}

// PublishEvent sends an event message.
//
// If ctx contains a message envelope, m is sent as a child of the message in
// that envelope.
func (s SinkSender) PublishEvent(
	ctx context.Context,
	m ax.Event,
	opts ...ax.PublishOption,
) (ax.Envelope, error) {
	env, err := s.newEnvelope(ctx, m)
	if err != nil {
		return ax.Envelope{}, err
	}

	for _, o := range opts {
		if err := o.ApplyPublishOption(&env); err != nil {
			return ax.Envelope{}, err
		}
	}

	return env, s.Sink.Accept(
		ctx,
		OutboundEnvelope{
			Envelope:  env,
			Operation: OpSendMulticast,
		},
	)
}

// newEnvelope returns an envelope containing m.
// The new envelope is configured as a child of the envelope in ctx, if any.
func (s SinkSender) newEnvelope(ctx context.Context, m ax.Message) (ax.Envelope, error) {
	validators := s.Validators
	if len(validators) == 0 {
		validators = DefaultValidators
	}

	for _, v := range validators {
		if err := v.Validate(ctx, m); err != nil {
			return ax.Envelope{}, err
		}
	}

	if env, ok := GetEnvelope(ctx); ok {
		return env.Envelope.NewChild(m), nil
	}

	return ax.NewEnvelope(m), nil
}
