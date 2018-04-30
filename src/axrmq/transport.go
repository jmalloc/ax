package axrmq

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
)

// Transport is an implementation of bus.Transport that uses RabbitMQ to
// communicate messages between endpoints.
type Transport struct {
}

// Initialize sets up the transport to communicate as an endpoint named ep.
func (t *Transport) Initialize(ctx context.Context, ep string) error {
	panic("not implemented")
}

// Subscribe instructs the transport to listen to multicast messages of the
// given type.
func (t *Transport) Subscribe(ctx context.Context, mt ax.MessageTypeSet) error {
	panic("not implemented")
}

// SendMessage sends a message.
func (t *Transport) SendMessage(ctx context.Context, m bus.OutboundEnvelope) error {
	panic("not implemented")
}

// ReceiveMessage returns the next message that has been delivered to the
// endpoint.
func (t *Transport) ReceiveMessage(ctx context.Context) (bus.InboundEnvelope, error) {
	panic("not implemented")
}
