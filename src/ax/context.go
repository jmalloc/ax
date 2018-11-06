package ax

import (
	"github.com/jmalloc/twelf/src/twelf"
)

// MessageContext provides context about the message being handled.
type MessageContext struct {
	// Envelope is the message envelope containing the message to be handled.
	Envelope Envelope

	// logger is used to log application-level messages about the handling of the
	// message.
	logger twelf.Logger
}

// NewMessageContext returns a message context for the given envelope.
func NewMessageContext(
	env Envelope,
	logger twelf.Logger,
) MessageContext {
	return MessageContext{
		env,
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
}
