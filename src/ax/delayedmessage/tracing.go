package delayedmessage

import (
	"time"

	"github.com/jmalloc/ax/src/ax"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// traceIntercept adds a log a event representing the fact that a message has
// been intercepted to be sent at a later tine.
func traceIntercept(span opentracing.Span, env ax.Envelope) {
	if span == nil {
		return
	}

	d := env.SendAt.Sub(env.CreatedAt)

	span.LogFields(
		log.String("event", "delayed-message.intercept"),
		log.String("message", "intercepting the message to be sent after a delay"),
		log.String("delay-for", d.String()),
		log.String("delay-until", env.SendAt.Format(time.RFC3339Nano)),
	)
}
