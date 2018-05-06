package bus

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// Transport is an interface for communicating messages between endpoints.
type Transport interface {
	MessageSender

	// Initialize sets up the transport to communicate as an endpoint named ep.
	Initialize(ctx context.Context, ep string) error

	// Subscribe instructs the transport to listen to multicast messages of the
	// given type.
	Subscribe(ctx context.Context, mt ax.MessageTypeSet) error

	// ReceiveMessage returns the next message that has been delivered to the
	// endpoint.
	ReceiveMessage(ctx context.Context) (InboundEnvelope, error)
}

// TransportStage is an outbound pipeline stage that forwards messages to a
// transport. It is typically used as the last stage in an outbound pipeline.
type TransportStage struct {
	transport Transport
}

// Initialize sets the transport used to send messages.
func (s *TransportStage) Initialize(ctx context.Context, t Transport) error {
	s.transport = t
	return nil
}

// SendMessage sends m via the transport.
func (s *TransportStage) SendMessage(ctx context.Context, m OutboundEnvelope) error {
	return s.transport.SendMessage(ctx, m)
}
