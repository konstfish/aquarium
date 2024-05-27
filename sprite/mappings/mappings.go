package mappings

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/aquarium/common/monitoring"
	"github.com/konstfish/aquarium/sprite/controllers"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var Router *gin.Engine

func CreateUrlMappings(rootDir string) {
	Router = gin.Default()
	Router.Use(otelgin.Middleware(monitoring.ServiceName, otelgin.WithFilter(monitoring.FilterTraces)))

	Router.Use(controllers.Cors())

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
