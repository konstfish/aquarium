package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/aquarium/common/db"
)

func Echo(c *gin.Context) {
	ctx := c.Request.Context()

	// send event to starfish
	db.Redis.PushToQueue(ctx, "starfish", "tetra")

	hostname, err := os.Hostname()
	if err != nil {
		c.String(http.StatusInternalServerError, "ERR")
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Hello! I'm %s", hostname))
}
