package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Healthz(c *gin.Context) {
	_, err := os.Hostname()
	if err != nil {
		c.String(http.StatusInternalServerError, "ERR")
		return
	}

	c.String(http.StatusOK, "OK")
}
