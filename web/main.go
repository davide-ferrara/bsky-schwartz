package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"web/internal/models"
	"web/views"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// Serve static files
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)

		// Get weights from session
		weights := make(map[string]float64)
		if w := session.Get("weights"); w != nil {
			if data, ok := w.([]byte); ok {
				json.Unmarshal(data, &weights)
			}
		}

		// Get flash message
		var message string
		if msg := session.Get("message"); msg != nil {
			message = msg.(string)
			session.Delete("message")
			session.Save()
		}

		component := views.SlidersPage("", message, weights)
		if err := component.Render(c.Request.Context(), c.Writer); err != nil {
			panic(err)
		}
	})

	r.GET("/values", func(c *gin.Context) {
		component := views.ValuesPage()
		if err := component.Render(c.Request.Context(), c.Writer); err != nil {
			panic(err)
		}
	})

	r.GET("/login", func(c *gin.Context) {
		component := views.LoginPage()
		if err := component.Render(c.Request.Context(), c.Writer); err != nil {
			panic(err)
		}
	})

	r.POST("/preferences", func(c *gin.Context) {
		session := sessions.Default(c)

		// Parse form data into map
		weights := make(map[string]float64)
		for _, value := range models.SwartzValues {
			if val := c.PostForm(value.ID); val != "" {
				var f float64
				if _, err := fmt.Sscanf(val, "%f", &f); err == nil {
					weights[value.ID] = f
				}
			}
		}

		// Save to session
		if data, err := json.Marshal(weights); err == nil {
			session.Set("weights", data)
			session.Set("message", "Preferenze salvate con successo!")
			session.Save()
		}

		c.Redirect(http.StatusFound, "/")
	})

	r.GET("/auth/bsky", func(c *gin.Context) {
		// TODO: Implement Bluesky AppPassword Auth
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
