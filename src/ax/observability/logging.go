package observability

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/twelf/src/twelf"
)

const (
	// identification fields
	messageIDIcon     = "=" // (equals) <message id>
	causationIDIcon   = "∵" // (because) <causation id>
	correlationIDIcon = "⋲" // (member of) <correlation id>

	// message logging
	inboundIcon       = "▼" // inbound messages can be considered as being (down)loaded
	inboundErrorIcon  = "▽"
	outboundIcon      = "▲" // outbound messages can be considered as being (up)loaded
	outboundErrorIcon = "△"
	retryIcon         = "↻" // shown when an inbound message is attempted for the 2nd+ time
	errorIcon         = "✖" // shown when handling or sending a message fails

	// application logging
	domainIcon     = "∴" // (therefore) to represent a decision made as a result of a message
	projectionIcon = "Σ" // (sum) to represent aggregating events

	// other
	// systemIcon    = "⚙" // (sprocket) to represent internals
	separatorIcon = "●" // a bold bullet used to separate arbitrary text fields
)

// LoggingObserver is an observer that logs about messages.
type LoggingObserver struct {
	logger twelf.Logger
}

// InitializeInbound initializes the observer for inbound messages.
func (o *LoggingObserver) InitializeInbound(ctx context.Context, ep *endpoint.Endpoint) error {
	o.logger = ep.Logger
	return nil
}

// InitializeOutbound initializes the observer for outbound messages.
func (o *LoggingObserver) InitializeOutbound(ctx context.Context, ep *endpoint.Endpoint) error {
	o.logger = ep.Logger
	return nil
}

// BeforeInbound logs information about an inbound message.
func (o *LoggingObserver) BeforeInbound(ctx context.Context, env endpoint.InboundEnvelope) {
	var stateIcon string
	if in, ok := endpoint.GetEnvelope(ctx); ok && in.AttemptCount > 1 {
		stateIcon = retryIcon
	}

	o.logMessage(inboundIcon, stateIcon, env.Envelope)
}

// AfterInbound logs information about errors that occur processing an inbound message.
func (o *LoggingObserver) AfterInbound(ctx context.Context, env endpoint.InboundEnvelope, err error) {
	if err != nil {
		o.logError(inboundErrorIcon, env.Envelope, err)
	}
}

// BeforeOutbound logs information about an outbound message.
func (o *LoggingObserver) BeforeOutbound(ctx context.Context, env endpoint.OutboundEnvelope) {
	o.logMessage(outboundIcon, "", env.Envelope)
}

// AfterOutbound logs information about an outbound message.
func (o *LoggingObserver) AfterOutbound(ctx context.Context, env endpoint.OutboundEnvelope, err error) {
	if err != nil {
		o.logError(outboundErrorIcon, env.Envelope, err)
	}
}

func (o *LoggingObserver) logMessage(typeIcon, statusIcon string, env ax.Envelope) {
	twelf.LogString(
		o.logger,
		formatLogMessage(
			typeIcon,
			statusIcon,
			env,
			env.Message.MessageDescription(),
		),
	)
}

func (o *LoggingObserver) logError(typeIcon string, env ax.Envelope, err error) {
	twelf.LogString(
		o.logger,
		formatLogMessage(
			typeIcon,
			errorIcon,
			env,
			err.Error(),
			env.Message.MessageDescription(),
		),
	)
}

// Logger is an implementation of twelf.Logger that adds information about the
// messaging being handled.
type Logger struct {
	Next     twelf.Logger
	Envelope ax.Envelope
	Icon     string
}

// NewDomainLogger returns a new logger that can be used to log information
// about the domain logic that results from handling env.Message.
func NewDomainLogger(
	next twelf.Logger,
	env ax.Envelope,
) twelf.Logger {
	return &Logger{
		next,
		env,
		domainIcon,
	}
}

// NewProjectionLogger returns a new logger that can be used to log information
// about the projections built as a result of applying env.Message.
func NewProjectionLogger(
	next twelf.Logger,
	env ax.Envelope,
) twelf.Logger {
	return &Logger{
		next,
		env,
		projectionIcon,
	}
}

// Log writes an application log message formatted according to a format
// specifier.
//
// It should be used for messages that are intended for people responsible for
// operating the application, such as the end-user or operations staff.
//
// f is the format specifier, as per fmt.Printf(), etc.
func (l *Logger) Log(f string, v ...interface{}) {
	l.LogString(fmt.Sprintf(f, v...))
}

// LogString writes a pre-formatted application log message.
//
// It should be ussed for messages that are intended for people responsible for
// operating the application, such as the end-user or operations staff.
func (l *Logger) LogString(s string) {
	twelf.Log(l.Next, l.format(s))
}

// Debug writes a debug log message formatted according to a format
// specifier.
//
// If IsDebug() returns false, no logging is performed.
//
// It should be used for messages that are intended for the software developers
// that maintain the application.
//
// f is the format specifier, as per fmt.Printf(), etc.
func (l *Logger) Debug(f string, v ...interface{}) {
	l.DebugString(fmt.Sprintf(f, v...))
}

// DebugString writes a pre-formatted debug log message.
//
// If IsDebug() returns false, no logging is performed.
//
// It should be used for messages that are intended for the software developers
// that maintain the application.
func (l *Logger) DebugString(s string) {
	twelf.DebugString(l.Next, l.format(s))
}

// IsDebug returns true if this logger will perform debug logging.
//
// Generally the application should just call Debug() or DebugString() without
// calling IsDebug(), however it can be used to check if debug logging is
// necessary before executing expensive code that is only used to obtain debug
// information.
func (l *Logger) IsDebug() bool {
	return twelf.IsDebug(l.Next)
}

func (l *Logger) format(s string) string {
	return formatLogMessage(
		l.Icon,
		"",
		l.Envelope,
		s,
	)
}

// formatLogMessage builds a log message that includes information about an Ax
// message.
func formatLogMessage(
	typeIcon string,
	statusIcon string,
	env ax.Envelope,
	messages ...string,
) string {
	var message string

	for _, m := range messages {
		message += " " + separatorIcon + " " + m
	}

	return fmt.Sprintf(
		"%s %s  %s %s  %s %s  %1s %1s  %s%s",
		messageIDIcon, env.MessageID,
		causationIDIcon, env.CausationID,
		correlationIDIcon, env.CorrelationID,
		typeIcon,
		statusIcon,
		env.Type(),
		message,
	)
}
