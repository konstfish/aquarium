package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Echo(c *gin.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error")
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Hello! I'm %s", hostname))
}
