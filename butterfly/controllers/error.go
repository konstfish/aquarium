package controllers

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context) {

	num := rand.Float64()

	if num < 0.5 {
		c.String(http.StatusServiceUnavailable, "This Request was not successful!")
		return
	}

	c.String(http.StatusOK, "This Request was successful!")
}
