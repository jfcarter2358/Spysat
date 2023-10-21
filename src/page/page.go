package page

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RedirectIndexPage(c *gin.Context) {
	c.Redirect(301, "/ui/home")
}

func ShowHomePage(c *gin.Context) {
	showPage(c, "home.html", gin.H{})
}

func ShowGroupPage(c *gin.Context) {
	showPage(c, "group.html", gin.H{})
}

func ShowObserverPage(c *gin.Context) {
	showPage(c, "group.html", gin.H{})
}

func showPage(c *gin.Context, page string, header gin.H) {
	render(c, header, page)
}

func render(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		c.XML(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}
