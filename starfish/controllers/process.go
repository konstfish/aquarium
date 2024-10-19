package controllers

import (
	"context"

	"github.com/konstfish/aquarium/common/logging"
	"github.com/konstfish/aquarium/common/monitoring"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.opentelemetry.io/otel/trace"
)

func ProcessEventAmount(ctx context.Context, event string) {
	var span trace.Span
	ctx, span = monitoring.Tracer.Start(ctx, "ProcessEventAmount")
	defer span.End()

	logging.Info(ctx, "Revieved event from", event)

	_ = ginmetrics.GetMonitor().GetMetric("starfish_event_amount").Inc([]string{event})
}
