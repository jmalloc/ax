package endpoint

import (
	"context"
	"time"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/internal/tracing"
	opentracing "github.com/opentracing/opentracing-go"
)

// InboundTransport is an interface receiving messages from endpoints.
type InboundTransport interface {
	// Initialize sets up the transport to communicate as an endpoint named ep.
	Initialize(ctx context.Context, ep string) error

	// Subscribe configures the transport to listen to messages of type mt that are
	// sent using op.
	Subscribe(ctx context.Context, op Operation, mt ax.MessageTypeSet) error

	// Receive returns the next message sent to this endpoint.
	// It blocks until a message is available, or ctx is canceled.
	Receive(ctx context.Context) (InboundEnvelope, Acknowledger, error)
}

// OutboundTransport is an interface for sending messages to endpoints.
type OutboundTransport interface {
	// Initialize sets up the transport to communicate as an endpoint named ep.
	Initialize(ctx context.Context, ep string) error

	// Send sends env via the transport.
	Send(ctx context.Context, env OutboundEnvelope) error
}

// TransportStage is an outbound pipeline stage that forwards messages to a
// transport. It is typically used as the last stage in an outbound pipeline.
type TransportStage struct {
	transport OutboundTransport
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (s *TransportStage) Initialize(ctx context.Context, ep *Endpoint) error {
	s.transport = ep.OutboundTransport
	return nil
}

// Accept sends env via the transport.
func (s *TransportStage) Accept(ctx context.Context, env OutboundEnvelope) error {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		tracing.LogEventS(
			span,
			"send",
			"sending the message via the transport",
			tracing.TypeName("pipeline_stage", s),
		)

		// if there is a span in the context, propagate it via the transport
		env.SpanContext = span.Context()
	}

	return s.transport.Send(ctx, env)
}

// Acknowledger is an interface for acknowledging a specific inbound message.
type Acknowledger interface {
	// Ack acknowledges the message, indicating that is was handled successfully
	// and does not need to be retried.
	Ack(ctx context.Context) error

	// Retry requeues the message so that it is retried at some point in the
	// future.
	//
	// d is a hint as to how long the transport should wait before retrying
	// this message.
	Retry(ctx context.Context, err error, d time.Duration) error

	// Reject indicates that the message could not be handled and should not be
	// retried. Depending on the transport, this may move the message to some form
	// of error queue or otherwise drop the message completely.
	Reject(ctx context.Context, err error) error
}
