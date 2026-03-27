package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"bsky-schwarz/pkg/bluesky"
	"bsky-schwarz/pkg/logger"
	"bsky-schwarz/pkg/middleware"
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

	logLevel := getEnvOrDefault("LOG_LEVEL", "info")
	if err := logger.Init(logLevel); err != nil {
		panic("failed to init logger: " + err.Error())
	}

	if err := scorer.Init(); err != nil {
		panic("failed to init scorer: " + err.Error())
	}

	logger.Info("server starting", "log_level", logLevel)

	handle := getEnv("BSKY_HANDLE")
	appPassword := getEnv("BSKY_APP_PASSWORD")

	bskyClient = bluesky.NewClient(handle, appPassword)

	router := gin.Default()
	router.Use(middleware.LoggingMiddleware())

	router.GET("/health", healthHandler)
	router.GET("/api/search", searchURIsHandler)
	router.GET("/api/analysis", analysisHandler)
	router.POST("/api/analysis/by-url", analysisByURLHandler)

	logger.Info("server listening", "port", 8080)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

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

	urls := bskyClient.GetPostsURL(query, limit)
	c.JSON(http.StatusOK, urls)
}

func analysisHandler(c *gin.Context) {
	start := time.Now()
	reqLogger := middleware.GetLogger(c)
	reqID := middleware.GetRequestID(c)

	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty query"})
		return
	}

	limitStr := c.DefaultQuery("limit", "1")

	modelKey := c.Query("model")
	if modelKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty model"})
		return
	}

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

	blueskyStart := time.Now()
	posts := bskyClient.QueryPosts(query, limit)
	blueskyDuration := time.Since(blueskyStart)

	reqLogger.Info("bluesky query timing",
		"request_id", reqID,
		"duration_ms", blueskyDuration.Milliseconds(),
		"posts_count", len(posts),
	)

	maxWorkers := scorer.GetConfig().Workers.MaxConcurrent
	if maxWorkers == 0 {
		maxWorkers = 10
	}

	numOfPosts := len(posts)
	if numOfPosts == 0 {
		c.JSON(http.StatusOK, []scorer.FeedItem{})
		return
	}

	reqLogger.Info("starting parallel analysis",
		"request_id", reqID,
		"workers", maxWorkers,
		"posts", numOfPosts,
		"model", model,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan scorer.FeedItem, numOfPosts)
	results := make(chan scorer.FeedItem, numOfPosts)
	errors := make(chan error, maxWorkers)

	for w := 0; w < maxWorkers; w++ {
		go func(workerID int) {
			for post := range jobs {
				postStart := time.Now()
				if err := post.ValueAlignment(model); err != nil {
					reqLogger.Error("worker analysis failed",
						"request_id", reqID,
						"worker_id", workerID,
						"post_uri", post.URI,
						"error", err,
						"duration_ms", time.Since(postStart).Milliseconds(),
					)
					errors <- fmt.Errorf("worker %d failed on post %s: %w", workerID, post.URI, err)
					return
				}
				reqLogger.Debug("worker analysis completed",
					"request_id", reqID,
					"worker_id", workerID,
					"post_uri", post.URI,
					"duration_ms", time.Since(postStart).Milliseconds(),
				)
				select {
				case results <- post:
				case <-ctx.Done():
					return
				}
			}
		}(w)
	}

	for j := 0; j < numOfPosts; j++ {
		jobs <- posts[j]
	}
	close(jobs)

	var res []scorer.FeedItem
	successCount := 0

	for i := 0; i < numOfPosts; i++ {
		select {
		case err := <-errors:
			cancel()
			reqLogger.Error("analysis failed, stopping all workers",
				"request_id", reqID,
				"error", err,
				"successful_posts", successCount,
				"total_duration_ms", time.Since(start).Milliseconds(),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Analysis failed: %v", err),
			})
			return
		case post := <-results:
			res = append(res, post)
			successCount++
		}
	}

	totalDuration := time.Since(start)
	reqLogger.Info("analysis request summary",
		"request_id", reqID,
		"query", query,
		"limit", limit,
		"model", model,
		"posts_found", len(posts),
		"posts_success", successCount,
		"bluesky_query_ms", blueskyDuration.Milliseconds(),
		"total_request_ms", totalDuration.Milliseconds(),
	)

	c.JSON(http.StatusOK, res)
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Batch analyze posts from Bluesky URLs
func analysisByURLHandler(c *gin.Context) {
	start := time.Now()
	reqLogger := middleware.GetLogger(c)
	reqID := middleware.GetRequestID(c)

	var req struct {
		URLs  []string `json:"urls"`
		Model string   `json:"model"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	if len(req.URLs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty urls list"})
		return
	}

	modelKey := req.Model
	if modelKey == "" {
		modelKey = "gpt"
	}

	model := scorer.GetConfig().Models[modelKey]
	if model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid model"})
		return
	}

	reqLogger.Info("starting batch analysis",
		"request_id", reqID,
		"urls_count", len(req.URLs),
		"model", model,
	)

	maxWorkers := scorer.GetConfig().Workers.MaxConcurrent
	if maxWorkers == 0 {
		maxWorkers = 10
	}

	numOfURLs := len(req.URLs)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan string, numOfURLs)
	results := make(chan *scorer.FeedItem, numOfURLs)
	errors := make(chan error, maxWorkers)

	for w := 0; w < maxWorkers; w++ {
		go func(workerID int) {
			for blueskyURL := range jobs {
				postStart := time.Now()
				post, err := bskyClient.GetPostFromBlueskyURL(blueskyURL)
				if err != nil {
					reqLogger.Error("failed to parse bluesky URL",
						"request_id", reqID,
						"worker_id", workerID,
						"url", blueskyURL,
						"error", err,
					)
					errors <- fmt.Errorf("worker %d failed to parse URL %s: %w", workerID, blueskyURL, err)
					return
				}

				if post == nil {
					errors <- fmt.Errorf("worker %d: post not found for URL %s", workerID, blueskyURL)
					return
				}

				if err := post.ValueAlignment(model); err != nil {
					reqLogger.Error("worker analysis failed",
						"request_id", reqID,
						"worker_id", workerID,
						"post_uri", post.URI,
						"error", err,
						"duration_ms", time.Since(postStart).Milliseconds(),
					)
					errors <- fmt.Errorf("worker %d failed on post %s: %w", workerID, post.URI, err)
					return
				}

				reqLogger.Debug("worker analysis completed",
					"request_id", reqID,
					"worker_id", workerID,
					"url", blueskyURL,
					"post_uri", post.URI,
					"duration_ms", time.Since(postStart).Milliseconds(),
				)

				select {
				case results <- post:
				case <-ctx.Done():
					return
				}
			}
		}(w)
	}

	for j := 0; j < numOfURLs; j++ {
		jobs <- req.URLs[j]
	}
	close(jobs)

	var res []*scorer.FeedItem
	successCount := 0

	for i := 0; i < numOfURLs; i++ {
		select {
		case err := <-errors:
			cancel()
			reqLogger.Error("batch analysis failed, stopping all workers",
				"request_id", reqID,
				"error", err,
				"successful_posts", successCount,
				"total_duration_ms", time.Since(start).Milliseconds(),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Analysis failed: %v", err),
			})
			return
		case post := <-results:
			res = append(res, post)
			successCount++
		}
	}

	totalDuration := time.Since(start)
	reqLogger.Info("batch analysis request summary",
		"request_id", reqID,
		"urls_count", len(req.URLs),
		"model", model,
		"posts_success", successCount,
		"total_request_ms", totalDuration.Milliseconds(),
	)

	c.JSON(http.StatusOK, res)
}

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

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
