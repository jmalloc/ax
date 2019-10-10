package axtest

import (
	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax"
	"github.com/jmalloc/ax/endpoint"
)

// ContainsMessage returns true if v contains a message equal to m.
func ContainsMessage(v []proto.Message, m proto.Message) bool {
	for _, x := range v {
		if proto.Equal(x, m) {
			return true
		}
	}

	return false
}

// ConsistsOfMessages returns true if a and b contain equal messages, regardless of order.
func ConsistsOfMessages(a []proto.Message, b ...proto.Message) bool {
	if len(a) != len(b) {
		return false
	}

	for _, m := range a {
		if !ContainsMessage(b, m) {
			return false
		}
	}

	return true
}

// EnvelopesEqual returns true if a and b are equivalent.
func EnvelopesEqual(a, b ax.Envelope) bool {
	if !a.CreatedAt.Equal(b.CreatedAt) {
		return false
	}

	if !a.SendAt.Equal(b.SendAt) {
		return false
	}

	if !proto.Equal(a.Message, b.Message) {
		return false
	}

	// ensure the "difficult to compare" values are equal so the remainder of
	// the struct can be compared using the equality operator.
	a.CreatedAt = b.CreatedAt
	a.SendAt = b.SendAt
	a.Message = b.Message

	return a == b
}

// ContainsEnvelope returns true if v contains an envelope equal to m.
func ContainsEnvelope(v []ax.Envelope, env ax.Envelope) bool {
	for _, x := range v {
		if EnvelopesEqual(x, env) {
			return true
		}
	}

	return false
}

// ConsistsOfEnvelopes returns true if a and b contain equal messages, regardless of order.
func ConsistsOfEnvelopes(a []ax.Envelope, b ...ax.Envelope) bool {
	if len(a) != len(b) {
		return false
	}

	for _, env := range a {
		if !ContainsEnvelope(b, env) {
			return false
		}
	}

	return true
}

// InboundEnvelopesEqual returns true if a and b are equivalent.
func InboundEnvelopesEqual(a, b endpoint.InboundEnvelope) bool {
	if !EnvelopesEqual(a.Envelope, b.Envelope) {
		return false
	}

	// ensure the "difficult to compare" values are equal so the remainder of
	// the struct can be compared using the equality operator.
	a.Envelope = b.Envelope

	return a == b
}

// ContainsInboundEnvelope returns true if v contains an envelope equal to m.
func ContainsInboundEnvelope(v []endpoint.InboundEnvelope, env endpoint.InboundEnvelope) bool {
	for _, x := range v {
		if InboundEnvelopesEqual(x, env) {
			return true
		}
	}

	return false
}

// ConsistsOfInboundEnvelopes returns true if a and b contain equal messages, regardless of order.
func ConsistsOfInboundEnvelopes(a []endpoint.InboundEnvelope, b ...endpoint.InboundEnvelope) bool {
	if len(a) != len(b) {
		return false
	}

	for _, env := range a {
		if !ContainsInboundEnvelope(b, env) {
			return false
		}
	}

	return true
}

// OutboundEnvelopesEqual returns true if a and b are equivalent.
func OutboundEnvelopesEqual(a, b endpoint.OutboundEnvelope) bool {
	if !EnvelopesEqual(a.Envelope, b.Envelope) {
		return false
	}

	// ensure the "difficult to compare" values are equal so the remainder of
	// the struct can be compared using the equality operator.
	a.Envelope = b.Envelope

	return a == b
}

// ContainsOutboundEnvelope returns true if v contains an envelope equal to m.
func ContainsOutboundEnvelope(v []endpoint.OutboundEnvelope, env endpoint.OutboundEnvelope) bool {
	for _, x := range v {
		if OutboundEnvelopesEqual(x, env) {
			return true
		}
	}

	return false
}

// ConsistsOfOutboundEnvelopes returns true if a and b contain equal messages, regardless of order.
func ConsistsOfOutboundEnvelopes(a []endpoint.OutboundEnvelope, b ...endpoint.OutboundEnvelope) bool {
	if len(a) != len(b) {
		return false
	}

	for _, env := range a {
		if !ContainsOutboundEnvelope(b, env) {
			return false
		}
	}

	return true
}
