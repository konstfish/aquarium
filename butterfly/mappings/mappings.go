package mappings

import (
	"github.com/gin-gonic/gin"
	"github.com/konstfish/aquarium/butterfly/controllers"
	"github.com/konstfish/aquarium/common/logging"
	"github.com/konstfish/aquarium/common/monitoring"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var Router *gin.Engine

func CreateUrlMappings() {
	Router = gin.New()
	Router.Use(gin.Recovery())

	// logging
	Router.Use(logging.CustomLogger())

	// opentelemetry
	Router.Use(otelgin.Middleware(monitoring.ServiceName, otelgin.WithFilter(monitoring.FilterTraces)))

	// metrics
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.Use(Router)

	// cors
	Router.Use(controllers.Cors())

	// health
	Router.GET("/healthz", controllers.Healthz)

	v1 := Router.Group("/butterfly/v1")
	{
		v1.POST("/error", controllers.Error)
	}
}
