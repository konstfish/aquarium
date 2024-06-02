package controllers

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/aquarium/common/db"
)

func Error(c *gin.Context) {
	ctx := c.Request.Context()

	// send event to starfish
	db.Redis.PushToQueue(ctx, "starfish", "butterfly")

	num := rand.Float64()

	if num < 0.5 {
		c.String(http.StatusServiceUnavailable, "This Request was not successful!")
		return
	}

	c.String(http.StatusOK, "This Request was successful!")
}
