package tracing

import (
	"context"
	"fmt"
	"reflect"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

// LogEvent logs an event to the span in ctx, if one is present.
func LogEvent(ctx context.Context, ev string, msg string, fields ...log.Field) {
	s := opentracing.SpanFromContext(ctx)
	LogEventS(s, ev, msg, fields...)
}

// LogEventS logs an event to span, if it is not nil.
func LogEventS(s opentracing.Span, ev string, msg string, fields ...log.Field) {
	if s == nil {
		return
	}

	fields = append(
		fields,
		log.String("event", ev),
		log.String("message", msg),
	)

	s.LogFields(fields...)
}

// LogError logs an error to the span in ctx, if one is present.
// The span is also marked as an error.
func LogError(ctx context.Context, err error, fields ...log.Field) {
	s := opentracing.SpanFromContext(ctx)
	LogErrorS(s, err, fields...)
}

// LogErrorS logs an error to span and marks it as an error, if it is not nil.
func LogErrorS(s opentracing.Span, err error, fields ...log.Field) {
	if s == nil {
		return
	}

	fields = append(
		fields,
		log.String("event", "error"),
		log.String("message", err.Error()),
		TypeName("error.kind", err),
	)

	ext.Error.Set(s, true)
	s.LogFields(fields...)
}

// TypeName returns a log field containing the name of v's type.
func TypeName(k string, v interface{}) log.Field {
	t := reflect.TypeOf(v)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return log.String(
		k,
		fmt.Sprintf(
			"%s.%s",
			t.PkgPath(),
			t.Name(),
		),
	)
}

// Duration returns a log field containing a human-readable representation of d.
func Duration(k string, d time.Duration) log.Field {
	return log.String(k, d.String())
}

// Time returns a log field containing a human-readable representation of t.
func Time(k string, t time.Time) log.Field {
	return log.String(k, t.Format(time.RFC3339Nano))
}
