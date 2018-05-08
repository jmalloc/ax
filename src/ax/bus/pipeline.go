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
	// Initialize is called after the transport is initialized. It can be used
	// to inspect or configure the transport as per the needs of the pipeline.
	Initialize(ctx context.Context, t Transport) error

	// Accept forwards an inbound message through the pipeline until
	// it is handled by some application-defined message handler(s).
	Accept(context.Context, MessageSink, InboundEnvelope) error
}

// OutboundPipeline is an interface for a message pipeline that processes
// messages that are sent via the message transport.
//
// A "stage" within the pipeline is simply an implementation of the
// OutboundPipeline interface that forwards messages to another pipeline.
type OutboundPipeline interface {
	MessageSink

	// Initialize is called after the transport is initialized. It can be used
	// to inspect or configure the transport as per the needs of the pipeline.
	Initialize(ctx context.Context, t Transport) error
}
