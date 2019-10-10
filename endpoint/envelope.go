package endpoint

import (
	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/ident"
	opentracing "github.com/opentracing/opentracing-go"
)

// AttemptID uniquely identifies an attempt to process a message.
type AttemptID struct {
	ident.ID
}

// GenerateAttemptID generates a new unique identifier for a processing attempt.
func GenerateAttemptID() AttemptID {
	var id AttemptID
	id.GenerateUUID()
	return id
}

// ParseAttemptID parses s into an attempt ID and returns it. It returns an
// error if s is empty.
func ParseAttemptID(s string) (AttemptID, error) {
	var id AttemptID
	err := id.Parse(s)
	return id, err
}

// MustParseAttemptID parses s into an attempt ID and returns it. It panics if
// s is empty.
func MustParseAttemptID(s string) AttemptID {
	var id AttemptID
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

	// AttemptID uniquely identifies the attempt to process this message.
	AttemptID AttemptID

	// AttemptCount is the number of times that an attempt has been made to process
	// this message.
	//
	// Messages may be retried after a failure handling the message, or if
	// an endpoint crashes, for example. Not all transports support an attempt
	// count. If the attempt count is unknown, it is set to zero.
	//
	// The attempt count may be reset if a message is manually re-queued after
	// being rejected by the retry policy.
	AttemptCount uint

	// SpanContext is the tracing context that was propagated with the message.
	SpanContext opentracing.SpanContext
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

	// SpanContext is the tracing context to propagate with the message.
	SpanContext opentracing.SpanContext
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
