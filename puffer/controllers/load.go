package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/aquarium/common/config"
	"github.com/konstfish/aquarium/common/db"
	"github.com/konstfish/aquarium/common/monitoring"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var BUTTERFLY_URL = config.GetConfigVar("BUTTERFLY_URL")
var TETRA_URL = config.GetConfigVar("TETRA_URL")

func Load(c *gin.Context) {
	ctx := c.Request.Context()

	var span trace.Span
	ctx, span = monitoring.Tracer.Start(ctx, "Load")
	defer span.End()

	lock := getLock(ctx)
	if lock != -1 {
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Next test can be sent in %ds", lock))
		return
	}

	setLock(ctx)

	err := callService(ctx, BUTTERFLY_URL, "POST", "callButterfly", 6)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Failed to ping butterfly")
		return
	}

	err = callService(ctx, TETRA_URL, "GET", "callTetra", 3)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Failed to ping tetra")
		return
	}

	c.String(http.StatusOK, "Pinged other services!")
}

func getLock(ctx context.Context) int {
	// get lock
	_, err := db.Redis.Client.Get(ctx, "load").Result()
	// if there is a lock, return 503
	if err == nil {
		// get remanining ttl
		ttl := db.Redis.Client.TTL(ctx, "load")
		seconds := ttl.Val().Abs().Seconds()

		return int(seconds)
	}

	return -1
}

func setLock(ctx context.Context) {
	num := fmt.Sprintf("%f", rand.Float64())

	db.Redis.Client.Set(ctx, "load", num, time.Second*30)
}

func callService(ctx context.Context, url string, method string, name string, num int) error {
	var span trace.Span
	ctx, span = monitoring.Tracer.Start(ctx, name)
	defer span.End()

	client := &http.Client{}

	for i := 0; i < num; i++ {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return err
		}

		// inject context
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

		_, err = client.Do(req)
		if err != nil {
			return err
		}
	}

	return nil
}
