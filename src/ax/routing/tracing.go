package routing

import (
	"github.com/jmalloc/ax/src/internal/reflectx"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func traceDispatch(span opentracing.Span, h MessageHandler) {
	if span == nil {
		return
	}

	span.LogFields(
		log.String("event", "routing.dispatch"),
		log.String("message", "dispatching message to handler"),
		log.String("handler", reflectx.PrettyTypeName(h)),
	)
}

func traceSelect(span opentracing.Span, dest string) {
	if span == nil {
		return
	}

	span.SetTag("message.destination", dest)

	span.LogFields(
		log.String("event", "routing.route"),
		log.String("message", "destination endpoint selected, forwarding message to the next pipeline stage"),
		log.String("endpoint", dest),
	)
}

func tracePreserve(span opentracing.Span, dest string) {
	if span == nil {
		return
	}

	span.SetTag("message.destination", dest)

	span.LogFields(
		log.String("event", "routing.route"),
		log.String("message", "destination endpoint already present in message, forwarding message to the next pipeline stage"),
		log.String("endpoint", dest),
	)
}

func traceForward(span opentracing.Span) {
	if span == nil {
		return
	}

	span.LogFields(
		log.String("event", "routing.route"),
		log.String("message", "message does not require a single destination, forwarding message to the next pipeline stage"),
	)
}
