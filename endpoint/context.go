package endpoint

import (
	"context"
)

// WithEnvelope returns a new context derived from p that contains env.
// The envelope can be retreived from the context with GetEnvelope().
func WithEnvelope(p context.Context, env InboundEnvelope) context.Context {
	return context.WithValue(
		p,
		envelopeKey,
		env,
	)
}

// GetEnvelope returns the message envelope contained in ctx.
// If ctx does not contain an envelope then ok is false.
func GetEnvelope(ctx context.Context) (env InboundEnvelope, ok bool) {
	v := ctx.Value(envelopeKey)

	if v != nil {
		env, ok = v.(InboundEnvelope)
	}

	return
}

// contextKey is a type used for the keys of context values. A specific type is
// used to prevent collisions with context keys from other packages.
type contextKey string

const (
	envelopeKey contextKey = "env"
)
