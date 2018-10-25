package endpoint

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/ident"
)

// DeliveryID uniquely identifies an attempt to process a message.
type DeliveryID struct {
	ident.ID
}

// GenerateDeliveryID generates a new unique identifier for a delivery.
func GenerateDeliveryID() DeliveryID {
	var id DeliveryID
	id.GenerateUUID()
	return id
}

// ParseDeliveryID parses s into a delivery ID and returns it. It returns an
// error if s is empty.
func ParseDeliveryID(s string) (DeliveryID, error) {
	var id DeliveryID
	err := id.Parse(s)
	return id, err
}

// MustParseDeliveryID parses s into a delivery ID and returns it. It panics if
// s is empty.
func MustParseDeliveryID(s string) DeliveryID {
	var id DeliveryID
	id.MustParse(s)
	return id
}

// InboundEnvelope is a specialization of ax.Envelope for messages that are
// received by an endpoint.
//
// Inbound envelopes traverse an InboundPipeline.
type InboundEnvelope struct {
	ax.Envelope

	// SourceEndpoint is the endpoint that sent the message.
	SourceEndpoint string

	// DeliveryID uniquely identifies the attempt to process this message.
	DeliveryID DeliveryID

	// DeliveryCount is the number of times that this message has been delivered.
	//
	// Messages may be redelivered after a failure handling the message, or if
	// an endpoint crashes, for example. Not all transports support a delivery
	// count, in which case the count is zero.
	//
	// The delivery count may be reset if a message is manually re-queued after
	// being rejected by the retry policy.
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
