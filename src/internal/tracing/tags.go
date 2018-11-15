package tracing

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
)

// SetTag sets a tag on the span in ctx, if one is present.
func SetTag(ctx context.Context, k string, v interface{}) {
	s := opentracing.SpanFromContext(ctx)
	SetTagS(s, k, v)
}

// SetTagS sets a tag on s, if it is not nil.
func SetTagS(s opentracing.Span, k string, v interface{}) {
	if s != nil {
		s.SetTag(k, v)
	}
}
