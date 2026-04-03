package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/log"
)

func NewServer(cfg *Config, logger *slog.Logger) http.Handler {
	// Crea mux per feed generator (http.ServeMux)
	feedMux := http.NewServeMux()

	// Feed generator endpoints
	feedMux.HandleFunc("GET /.well-known/did.json", handleWellKnown(cfg))
	feedMux.HandleFunc("GET /xrpc/_health", handleHealth)
	feedMux.HandleFunc("GET /xrpc/app.bsky.feed.describeFeedGenerator", handleDescribeFeedGenerator(cfg))
	feedMux.HandleFunc("GET /xrpc/app.bsky.feed.getFeedSkeleton", handleGetFeedSkeleton(cfg, logger))

	// Crea router Gin per web app
	gin.SetMode(gin.ReleaseMode)
	webRouter := gin.New()
	webRouter.Use(gin.Recovery())

	// Session middleware
	sessionSecret := os.Getenv("SESSION_KEY")
	if sessionSecret == "" {
		err := fmt.Errorf("you must coonfigure SESSION_KEY in .env")
		log.Error(err)
		os.Exit(1)
	}

	store := cookie.NewStore([]byte(sessionSecret))
	webRouter.Use(sessions.Sessions("bluesky-session", store))

	// Static files
	webRouter.Static("/static", "./static")
	webRouter.Static("/lexicons", "./lexicons")

	// Web app routes
	webRouter.GET("/", RootHandler)
	webRouter.GET("/login", LoginGetHandler)
	webRouter.POST("/login", LoginPostHandler)
	webRouter.POST("/logout", LogoutHandler)
	webRouter.GET("/profile", ProfileHandler)
	webRouter.POST("/delete-weights", DeleteWeightsHanlder)
	webRouter.GET("/values", ValuesHandler)
	webRouter.POST("/preferences", PreferencesHandler)

	// Route specifiche del feed generator
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/.well-known/did.json" ||
			r.URL.Path == "/xrpc/_health" ||
			r.URL.Path == "/xrpc/app.bsky.feed.describeFeedGenerator" ||
			r.URL.Path == "/xrpc/app.bsky.feed.getFeedSkeleton" {
			feedMux.ServeHTTP(w, r)
			return
		}

		// Tutto il resto va alla web app
		webRouter.ServeHTTP(w, r)
	})
}
