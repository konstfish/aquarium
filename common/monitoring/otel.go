package monitoring

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/konstfish/aquarium/common/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer
var ServiceName string

func InitTracer(service string) {
	ServiceName = service
	Tracer = otel.Tracer(service)

	if config.GetConfigVar("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
		return
	}

	ctx := context.Background()

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(config.GetConfigVar("OTEL_EXPORTER_OTLP_ENDPOINT")),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
		)),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	log.Println("Tracing initialized")
}

func ExtractTraceparentHeader(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()

	if !sc.TraceID().IsValid() || !sc.SpanID().IsValid() {
		return ""
	}

	return fmt.Sprintf("00-%s-%s-%s", sc.TraceID(), sc.SpanID(), sc.TraceFlags())
}

func EmptyTraceparentHeader() string {
	return "00-00000000000000000000000000000000-0000000000000000-00"
}

func ParseTraceparentHeader(traceparentHeader string) (trace.SpanContext, error) {
	parts := strings.Split(traceparentHeader, "-")
	if len(parts) != 4 {
		return trace.SpanContext{}, fmt.Errorf("invalid traceparent header")
	}

	traceID, err := trace.TraceIDFromHex(parts[1])
	if err != nil {
		return trace.SpanContext{}, fmt.Errorf("invalid TraceID: %w", err)
	}

	spanID, err := trace.SpanIDFromHex(parts[2])
	if err != nil {
		return trace.SpanContext{}, fmt.Errorf("invalid SpanID: %w", err)
	}

	traceFlagsInt, err := strconv.ParseUint(parts[3], 16, 8)
	if err != nil {
		return trace.SpanContext{}, fmt.Errorf("invalid traceFlags: %w", err)
	}

	traceFlags := byte(traceFlagsInt)

	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.TraceFlags(traceFlags),
		Remote:     true,
	}), nil
}

// middleware
func FilterTraces(req *http.Request) bool {
	// check if request method is options
	if req.Method == "OPTIONS" {
		return false
	}

	// check if request path is noisy
	var notToLogEndpoints = []string{"/healthz", "/metrics"}
	if slices.Index(notToLogEndpoints, req.URL.Path) != -1 {
		return false
	}

	return true
}
