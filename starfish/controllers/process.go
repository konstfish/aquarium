package controllers

import (
	"context"
	"fmt"
	"log"

	"github.com/konstfish/aquarium/common/monitoring"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.opentelemetry.io/otel/trace"
)

func ProcessEventAmount(ctx context.Context, event string) {
	var span trace.Span
	ctx, span = monitoring.Tracer.Start(ctx, "ProcessEventAmount")
	defer span.End()

	log.Println(fmt.Sprintf("Revieved event from %s", event))

	_ = ginmetrics.GetMonitor().GetMetric("starfish_event_amount").Inc([]string{event})
}
