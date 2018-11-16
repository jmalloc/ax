package tracing

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// StartChildOf starts a new span as a child of the span in ctx, if any.
func StartChildOf(
	ctx context.Context,
	tr opentracing.Tracer,
	op string,
	opts ...opentracing.StartSpanOption,
) opentracing.Span {
	if p := opentracing.SpanFromContext(ctx); p != nil {
		if tr == nil {
			tr = p.Tracer()
		}

		opts = append(
			opts,
			opentracing.ChildOf(p.Context()),
		)
	}

	return StartSpan(tr, op, opts...)
}

// StartFollowsFrom starts a new span that follows from the span in ctx, if any.
func StartFollowsFrom(
	ctx context.Context,
	tr opentracing.Tracer,
	op string,
	opts ...opentracing.StartSpanOption,
) opentracing.Span {
	if p := opentracing.SpanFromContext(ctx); p != nil {
		if tr == nil {
			tr = p.Tracer()
		}

		opts = append(
			opts,
			opentracing.FollowsFrom(p.Context()),
		)
	}

	return StartSpan(tr, op, opts...)
}

// StartSpan starts a new span using tr, or a NoopTracer if tr is nil.
func StartSpan(
	tr opentracing.Tracer,
	op string,
	opts ...opentracing.StartSpanOption,
) opentracing.Span {
	if tr == nil {
		tr = opentracing.NoopTracer{}
	}

	opts = append(
		opts,
		opentracing.Tag{
			Key:   string(ext.Component),
			Value: "ax",
		},
	)

	return tr.StartSpan(op, opts...)
}
