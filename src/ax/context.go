package ax

import (
	"fmt"

	"github.com/jmalloc/ax/src/internal/tracing"
	"github.com/jmalloc/twelf/src/twelf"
	opentracing "github.com/opentracing/opentracing-go"
)

// MessageContext provides context about the message being handled.
type MessageContext struct {
	// Envelope is the message envelope containing the message to be handled.
	Envelope Envelope

	// span is the opentracing span for the current message.
	span opentracing.Span

	// logger is used to log application-level messages about the handling of the
	// message.
	logger twelf.Logger
}

// NewMessageContext returns a message context for the given envelope.
func NewMessageContext(
	env Envelope,
	span opentracing.Span,
	logger twelf.Logger,
) MessageContext {
	return MessageContext{
		env,
		span,
		logger,
	}
}

// Log writes an application-level log message about the handling of the
// message.
//
// The log messages should be understood by non-developers who are familiar with
// the application's business domain.
func (c *MessageContext) Log(f string, v ...interface{}) {
	twelf.Log(c.logger, f, v...)

	tracing.LogEventS(
		c.span,
		"application_log",
		fmt.Sprintf(f, v...),
	)
}
