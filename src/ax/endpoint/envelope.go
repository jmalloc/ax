package endpoint

import (
	"github.com/jmalloc/ax/src/ax"
)

// InboundEnvelope is a specialization of ax.Envelope for messages that are
// received by an endpoint.
//
// Inbound envelopes traverse an InboundPipeline.
type InboundEnvelope struct {
	ax.Envelope

	// SourceEndpoint is the endpoint that sent the message.
	SourceEndpoint string

	// DeliveryCount is the number of times that this message has been delivered.
	//
	// Messages may be redelivered after a failure handling the message, or if
	// an endpoint crashes, for example. Not all transports support a delivery
	// count, in which case the count is zero.
	DeliveryCount uint
}

// OutboundEnvelope is a specialization of ax.Envelope for messages that are
// sent by an endpoint.
//
// Outbound envelopes traverse an OutboundPipeline.
type OutboundEnvelope struct {
	ax.Envelope

	// Operation is the operation to be performed on the message. It dictates
	// how the message is sent by the transport.
	Operation Operation

	// DestinationEndpoint is the endpoint to which the message is sent when
	// Operation is OpSendUnicast. The field is ignored for other operations.
	DestinationEndpoint string
}

// Operation is an enumeration of transport operations that can be performed
// in order to send an outbound message.
type Operation int

const (
	// OpSendUnicast is an outbound transport operation that sends a message to
	// a specific endpoint as determined by the outbound message's
	// DestinationEndpoint property.
	OpSendUnicast Operation = iota

	// OpSendMulticast is an outbound transport operation that sends a message
	// to all of its subscribers.
	OpSendMulticast
)
