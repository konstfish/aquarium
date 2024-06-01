package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/aquarium/common/db"
)

func Healthz(c *gin.Context) {
	err := db.Redis.Client.Ping(c.Request.Context()).Err()
	if err != nil {
		c.String(http.StatusInternalServerError, "ERR")
		return
	}

	c.String(http.StatusOK, "OK")
}
