package mappings

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/aquarium/sprite/controllers"
)

var Router *gin.Engine

func CreateUrlMappings(rootDir string) {
	Router = gin.Default()
	// Router.Use(otelgin.Middleware(monitoring.ServiceName, otelgin.WithFilter(monitoring.FilterTraces)))

	Router.Use(controllers.Cors())

	controllers.InitSprites(rootDir)

	v1 := Router.Group("/sprites/v1")
	{
		for folder, files := range controllers.Sprites {
			route := "/" + folder
			v1.GET(route, func(c *gin.Context) {

				c.String(http.StatusOK, controllers.GetRandomSprite(rootDir+"/"+folder, files))
			})
		}
	}
}
