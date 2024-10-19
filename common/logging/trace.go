package logging

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func extractTrace(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()

	return traceID
}
