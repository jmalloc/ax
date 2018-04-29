package bus

import (
	"context"

	"github.com/jmalloc/ax/src/ax"
)

// InboundEnvelope is a specialization of ax.Envelope for messages that are
// received by "this" endpoint.
type InboundEnvelope struct {
	ax.Envelope

	// DeliveryCount is the number of times that this message has been delivered
	// to the endpoint.
	//
	// Messages may be redelivered after a failure handling the message, or of
	// an endpoint crashes, for example. Not all transports support a delivery
	// count, in which case the count is zero.
	DeliveryCount uint

	// Done is called to indicate that the endpoint has finished processing the
	// message. This does not necessarily mean that the message has been handled
	// successfully. The InboundOperation passed to Done() determines whether or
	// not the message is retried or not.
	Done func(context.Context, InboundOperation) error
}

// InboundOperation is an enumeration of operations that can be performed to
// an inbound message.
type InboundOperation int

const (
	// OpAck is an inbound transport operation that causes the inbound message
	// to be removed from the endpoint's queue.
	OpAck InboundOperation = iota

	// OpRetry is an inbound transport operation that causes the inbound message
	// to be retried.
	OpRetry

	// OpReject is an inbound transport operation that causes the inbound
	// message to be rejected. Depending on the transport, the message may be
	// moved to some form of error queue, or dropped completely.
	OpReject
)

// OutboundEnvelope is a specialization of ax.Envelope for messages that are
// sent by "this" endpoint.
type OutboundEnvelope struct {
	ax.Envelope

	// Operation is the operation to be performed on the message. It dictates
	// how the message is sent by the transport.
	Operation OutboundOperation

	// DestinationEndpoint is the endpoint to which the message is sent when
	// Operation is OpSendUnicast. The field is ignored for other operations.
	DestinationEndpoint string
}

// OutboundOperation is an enumeration of operations that can be performed to
// an outbound message.
type OutboundOperation int

const (
	// OpSendUnicast is an outbound transport operation that sends a message to
	// a specific endpoint as determined by the outbound message's
	// DestinationEndpoint property.
	OpSendUnicast OutboundOperation = iota

	// OpSendMulticast is an outbound transport operation that sends a message
	// to all of its subscribers.
	OpSendMulticast
)
