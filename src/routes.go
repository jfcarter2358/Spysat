// routes.go

package main

import (
	"net/http"
	"spysat/api"
	"spysat/middleware"
	"spysat/page"

	"github.com/gin-gonic/gin"
)

func initializeRoutes() {
	router.Static("/static/css", "./static/css")
	router.Static("/static/img", "./static/img")
	router.Static("/static/js", "./static/js")

	router.GET("/", page.RedirectIndexPage)

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
	})

	healthRoutes := router.Group("/health", middleware.CORSMiddleware())
	{
		healthRoutes.GET("/healthy", api.GetHealth)
	}

	apiRoutes := router.Group("/api", middleware.CORSMiddleware())
	{
		v1Routes := apiRoutes.Group("/v1")
		{
			baseStationRoutes := v1Routes.Group("/basestation")
			{
				baseStationRoutes.POST("/:group/:observer/:stream", api.UpdateBaseStation)
				baseStationRoutes.GET("/:group/:observer/:stream", api.GetBaseStation)
			}
			operationRoutes := v1Routes.Group("/operation")
			{
				operationRoutes.GET("/operation", api.GetOperation)
				operationRoutes.POST("/operation", api.UpdateOperation)
			}
			v1Routes.GET("/status", api.GetStatus)
		}
	}

	uiRoutes := router.Group("/ui", middleware.CORSMiddleware())
	{
		uiRoutes.GET("/", page.ShowHomePage)
		uiRoutes.GET("/:group", page.ShowGroupPage)
		uiRoutes.GET("/:group/:observer", page.ShowObserverPage)
	}
}
