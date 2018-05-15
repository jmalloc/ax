package endpoint

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Transport is an interface for communicating messages between endpoints.
type Transport interface {
	// Initialize sets up the transport to communicate as an endpoint named ep.
	Initialize(ctx context.Context, ep string) error

	// Subscribe configures the transport to listen to messages of type mt that are
	// sent using op.
	Subscribe(ctx context.Context, op Operation, mt ax.MessageTypeSet) error

	// Send sends env via the transport.
	Send(ctx context.Context, env OutboundEnvelope) error

	// Receive returns the next message sent to this endpoint.
	// It blocks until a message is available, or ctx is canceled.
	Receive(ctx context.Context) (InboundEnvelope, Acknowledger, error)
}

// TransportStage is an outbound pipeline stage that forwards messages to a
// transport. It is typically used as the last stage in an outbound pipeline.
type TransportStage struct {
	transport Transport
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (s *TransportStage) Initialize(ctx context.Context, ep *Endpoint) error {
	s.transport = ep.Transport
	return nil
}

// Accept sends env via the transport.
func (s *TransportStage) Accept(ctx context.Context, env OutboundEnvelope) error {
	return s.transport.Send(ctx, env)
}

// Acknowledger is an interface for acknowledging a specific inbound message.
type Acknowledger interface {
	// Ack acknowledges the message, indicating that is was handled successfully
	// and does not need to be redelivered.
	Ack(ctx context.Context) error

	// Retry requeues the message so that it is redelivered at some point in the
	// future.
	Retry(ctx context.Context, err error) error

	// Reject indicates that the message could not be handled and should not be
	// retried. Depending on the transport, this may move the message to some form
	// of error queue or otherwise drop the message completely.
	Reject(ctx context.Context, err error) error
}
