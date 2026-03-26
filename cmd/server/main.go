package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"bsky-schwarz/pkg/bluesky"
	"bsky-schwarz/pkg/scorer"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var bskyClient bluesky.Client

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("failed to load .env: %w", err))
	}

	if err := scorer.Init(); err != nil {
		panic("failed to init scorer: " + err.Error())
	}

	handle := getEnv("BSKY_HANDLE")
	appPassword := getEnv("BSKY_APP_PASSWORD")

	bskyClient = bluesky.NewClient(handle, appPassword)

	router := gin.Default()

	router.GET("/health", healthHandler)
	router.GET("/api/search", searchURIsHandler)
	router.GET("/api/analysis", analysisHandler)
	router.GET("/api/analysis/by-uri", analysisByUriHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

// Get a list of uri given specific query
func searchURIsHandler(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty query"})
		return
	}

	limitStr := c.DefaultQuery("limit", "1")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	feed := bskyClient.GetPostsUri(query, limit)
	c.JSON(http.StatusOK, feed)
}

// Get feed with values
func analysisHandler(c *gin.Context) {
	query := c.DefaultQuery("query", "test")
	limitStr := c.DefaultQuery("limit", "1")
	modelKey := c.DefaultQuery("model", "gpt")

	model := scorer.GetConfig().Models[modelKey]
	if model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid model"})
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	posts := bskyClient.QueryPosts(query, limit)

	for i := range posts {
		if err := posts[i].ValueAlignment(model); err != nil {
			log.Printf("Error analyzing post %s: %v", posts[i].URI, err)
		}
	}

	c.JSON(http.StatusOK, posts)
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Get post with values from URI
func analysisByUriHandler(c *gin.Context) {
	uri := c.Query("uri")
	if uri == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty uri"})
		return
	}

	modelKey := c.DefaultQuery("model", "gpt")
	model := scorer.GetConfig().Models[modelKey]
	if model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid model"})
		return
	}

	post := bskyClient.GetPost(uri)
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	if err := post.ValueAlignment(model); err != nil {
		log.Printf("Error analyzing post %s: %v", post.URI, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to analyze post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func getEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	panic("missing environment variable: " + key)
}
