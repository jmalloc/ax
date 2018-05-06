package bus

import (
	"context"
)

// InboundPipeline is an interface for a message pipeline that processes
// messages received from the message transport.
//
// A "stage" within the pipeline is simply an implementation of the
// InboundPipeline interface that forwards messages to another pipeline.
type InboundPipeline interface {
	// Initialize is called during endpoint initialization, after the transport
	// has been initialized.
	Initialize(ctx context.Context, t Transport) error

	// DeliverMessage forwards an inbound message through the pipeline until
	// it is handled by some application-defined message handler(s).
	DeliverMessage(ctx context.Context, s MessageSender, m InboundEnvelope) error
}

// OutboundPipeline is an interface for a message pipeline that processes
// messages that are sent via the message transport.
//
// A "stage" within the pipeline is simply an implementation of the
// OutboundPipeline interface that forwards messages to another pipeline.
type OutboundPipeline interface {
	MessageSender

	// Initialize is called during endpoint initialization, after the transport
	// has been initialized.
	Initialize(ctx context.Context, t Transport) error
}
