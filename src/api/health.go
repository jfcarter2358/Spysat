package api

import (
	"net/http"
	"spysat/health"

	"github.com/gin-gonic/gin"
)

func GetHealth(ctx *gin.Context) {
	if health.IsHealthy {
		ctx.Status(http.StatusOK)
	} else {
		ctx.Status(http.StatusInternalServerError)
	}
	return
}
