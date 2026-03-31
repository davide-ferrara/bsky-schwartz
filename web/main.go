package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"web/internal/models"
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

	r.POST("/preferences", func(c *gin.Context) {
		// Parse form data
		formData := make(map[string]string)
		for _, value := range models.SwartzValues {
			formData[value.ID] = c.PostForm(value.ID)
		}

		// Map to SchwartzWeights
		weights := models.MapFormToWeights(formData)

		// Log for now (TODO: save to database)
		fmt.Printf("Received preferences: %+v\n", weights)

		// Return success response
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Preferenze salvate con successo",
			"weights": weights,
		})
	})

	r.GET("/auth/bsky", func(c *gin.Context) {
		// TODO: Implement Bluesky OAuth
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	r.Run(":8080")
}
