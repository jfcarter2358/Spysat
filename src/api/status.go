package api

import (
	"net/http"
	"spysat/operation"

	"github.com/gin-gonic/gin"
)

func GetStatus(ctx *gin.Context) {
	out := make(map[string]map[string]string)
	for obName, ob := range operation.Operations.Observers {
		out[obName] = make(map[string]string)
		for sName, s := range ob.Streams {
			out[obName][sName] = s.Status
		}
	}

	ctx.JSON(http.StatusOK, out)
}
