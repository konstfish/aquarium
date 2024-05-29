package mappings

import (
	"github.com/gin-gonic/gin"
	"github.com/konstfish/aquarium/common/monitoring"
	"github.com/konstfish/aquarium/tetra/controllers"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var Router *gin.Engine

func CreateUrlMappings() {
	Router = gin.Default()
	Router.Use(otelgin.Middleware(monitoring.ServiceName, otelgin.WithFilter(monitoring.FilterTraces)))

	Router.Use(controllers.Cors())

	v1 := Router.Group("/tetra/v1")
	{
		v1.GET("/echo", controllers.Echo)
	}
}
