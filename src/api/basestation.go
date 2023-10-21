package api

import (
	"net/http"
	"spysat/basestation"
	"spysat/utils"

	"github.com/gin-gonic/gin"
)

func UpdateBaseStation(ctx *gin.Context) {
	group := ctx.Param("group")
	observer := ctx.Param("observer")
	name := ctx.Param("stream")

	var data map[string]string
	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.Error(err, ctx, http.StatusInternalServerError)
		return
	}

	if err := basestation.ReceiveData(data["data"], name, observer, group); err != nil {
		utils.Error(err, ctx, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GetBaseStation(ctx *gin.Context) {
	group := ctx.Param("group")
	observer := ctx.Param("observer")
	name := ctx.Param("stream")

	data := basestation.GetData(name, observer, group)

	ctx.JSON(http.StatusOK, data)
}
