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
	router.GET("/search/post", searchHandler)

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func searchHandler(c *gin.Context) {
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

func getEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	panic("missing environment variable: " + key)
}
