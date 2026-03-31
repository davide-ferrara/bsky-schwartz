package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"web/views"
)

func main() {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		component := views.SlidersPage("")
		component.Render(c.Request.Context(), c.Writer)
	})

	r.GET("/values", func(c *gin.Context) {
		component := views.ValuesPage()
		component.Render(c.Request.Context(), c.Writer)
	})

	r.GET("/login", func(c *gin.Context) {
		component := views.LoginPage()
		component.Render(c.Request.Context(), c.Writer)
	})

	r.GET("/auth/bsky", func(c *gin.Context) {
		// TODO: Implement Bluesky OAuth
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	r.Run(":8080")
}
