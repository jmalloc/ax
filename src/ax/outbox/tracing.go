package outbox

import (
	"github.com/jmalloc/ax/src/ax/endpoint"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// traceFound adds a log a event representing the fact that an outbox has been
// found for the inbound message.
func traceFound(span opentracing.Span, envs []endpoint.OutboundEnvelope) {
	n := len(envs)

	span.LogFields(
		log.String("event", "outbox.found"),
		log.String("message", "found an outbox, the message has already been processed"),
		log.Int("message-count", n),
	)
}

// traceNotFound adds a log a event representing the fact that no outbox has
// been found for the inbound message.
func traceNotFound(span opentracing.Span) {
	span.LogFields(
		log.String("event", "outbox.not-found"),
		log.String("message", "no outbox found, forwarding message to the next pipeline stage"),
	)
}

// traceSave adds a log event representing the fact that an outbox has been
// saved for the inbound message.
func traceSave(span opentracing.Span, envs []endpoint.OutboundEnvelope) {
	n := len(envs)

	span.LogFields(
		log.String("event", "outbox.save"),
		log.String("message", "saving an outbox for this message"),
		log.Int("message-count", n),
	)
}
