package mappings

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/aquarium/common/logging"
	"github.com/konstfish/aquarium/common/monitoring"
	"github.com/konstfish/aquarium/sprite/controllers"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var Router *gin.Engine

func CreateUrlMappings(rootDir string) {
	Router = gin.New()
	Router.Use(gin.Recovery())

	// opentelemetry
	Router.Use(otelgin.Middleware(monitoring.ServiceName, otelgin.WithFilter(monitoring.FilterTraces)))

	// metrics
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.Use(Router)

	// logging, needs to be after otel middleware to make sure we have a trace id in the context
	Router.Use(logging.CustomLogger())

	// cors
	Router.Use(controllers.Cors())

	// health
	Router.GET("/healthz", controllers.Healthz)

	controllers.InitSprites(rootDir)

	v1 := Router.Group("/sprites/v1")
	{
		for folder, files := range controllers.Sprites {
			route := "/" + folder

			v1.GET(route, func(c *gin.Context) {
				ctx := c.Request.Context()

				c.String(http.StatusOK, controllers.GetRandomSprite(ctx, rootDir+"/"+folder, files))
			})
		}
	}
}
