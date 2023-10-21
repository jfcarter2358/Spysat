package api

import (
	"net/http"
	"spysat/operation"
	"spysat/utils"

	"github.com/gin-gonic/gin"
)

func UpdateOperation(ctx *gin.Context) {
	var data operation.Operation
	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.Error(err, ctx, http.StatusInternalServerError)
		return
	}

	operation.Operations = data
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GetOperation(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, operation.Operations)
}
