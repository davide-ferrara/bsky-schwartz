package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type ModelConfig struct {
	Name     string
	Provider string // "openrouter" or "siliconflow"
	ModelID  string // model ID for the provider
}

var modelConfigs = []ModelConfig{
	{Name: "gpt-4.1-mini", Provider: "openrouter", ModelID: "openai/gpt-4.1-mini"},
	{Name: "mistral-14b", Provider: "openrouter", ModelID: "mistralai/ministral-14b-2512"},
	{Name: "deepseek-v3.2", Provider: "siliconflow", ModelID: "deepseek-ai/DeepSeek-V3"},
	{Name: "qwen3-vl-30b", Provider: "siliconflow", ModelID: "Qwen/Qwen3-VL-30B-A3B-Instruct"},
}

var modelFlag = flag.String("model", "", "Run analysis for specific model (partial name match)")
var limitFlag = flag.Int("limit", 2, "Number of posts to analyze per model")
var feedFlag = flag.String("feed", "./feed_20260330212226.json", "Path to feed JSON file")

func filterModels(configs []ModelConfig, name string) []ModelConfig {
	name = strings.ToLower(name)
	var filtered []ModelConfig
	for _, cfg := range configs {
		if strings.Contains(strings.ToLower(cfg.Name), name) {
			filtered = append(filtered, cfg)
		}
	}
	return filtered
}

func main() {
	flag.Parse()

	if err := godotenv.Load(".env"); err != nil {
		fmt.Fprintf(os.Stderr, "warning: .env not found: %v\n", err)
	}

	configs := modelConfigs
	if *modelFlag != "" {
		configs = filterModels(modelConfigs, *modelFlag)
		if len(configs) == 0 {
			fmt.Printf("No models matching '%s' found\n", *modelFlag)
			return
		}
		fmt.Printf("Running analysis for %d model(s) matching '%s'\n", len(configs), *modelFlag)
	}

	posts, err := LoadStaticPosts(*feedFlag)
	if err != nil {
		fmt.Println("Could not open JSON")
		panic(err)
	}

	fmt.Println("========================================")
	fmt.Println("Starting analysis...")
	fmt.Printf("Loaded %d posts from %s\n", len(posts), *feedFlag)
	fmt.Printf("Analyzing %d posts per model\n", *limitFlag)
	fmt.Println("========================================")

	for i, cfg := range configs {
		fmt.Printf("\n[%d/%d] Model: %s (%s)\n", i+1, len(configs), cfg.Name, cfg.Provider)
		fmt.Println("----------------------------------------")

		var client AIClient
		var err error

		switch cfg.Provider {
		case "openrouter":
			client, err = GetOpenRouterClient()
		case "siliconflow":
			client, err = GetSiliconFlowClient()
		default:
			fmt.Printf("ERROR: Unknown provider: %s\n", cfg.Provider)
			continue
		}

		if err != nil {
			fmt.Printf("ERROR: Failed to init client: %v\n", err)
			continue
		}

		parts := strings.Split(cfg.ModelID, "/")
		filename := fmt.Sprintf("post_%s", parts[len(parts)-1])
		if err := runAnalysisSync(client, posts[:*limitFlag], cfg.ModelID, filename); err != nil {
			fmt.Printf("ERROR: Analysis failed: %v\n", err)
		}
	}

	fmt.Println("\n========================================")
	fmt.Println("All models processed.")
	fmt.Println("========================================")
}

func DownloadFeed(path string) []Post {
	ctx := context.Background()
	handle := GetEnv("BSKY_HANDLE")
	appPassword := GetEnv("BSKY_APP_PASSWORD")
	bskyClient, err := NewClient(handle, appPassword)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting to Download Feed...")

	postUrls, err := LoadPostURLs(path)
	if err != nil {
		panic(err)
	}

	var posts []Post
	for _, postURL := range postUrls.Urls {
		post, err := bskyClient.GetPostUrl(ctx, postURL)
		if err != nil {
			fmt.Errorf("Error: %w", err)
			continue
		}

		fmt.Println("Appending post:", post.AtURI)

		posts = append(posts, post)
	}

	filename := strings.Split(path, "/")[1]
	SavePostsToJson(filename, posts)

	fmt.Println("Saved file to:", filename)

	return posts
}

func runAnalysisSync(c AIClient, posts []Post, model string, filename string) error {
	startTime := time.Now()

	for i := range posts {
		postStart := time.Now()
		fmt.Printf("  [Post %d/%d] Analyzing: %s\n", i+1, len(posts), truncate(posts[i].Text, 50))

		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)

		analysis, err := CalculateRating(ctx, c, model, &posts[i])
		if err != nil {
			posts[i].ValueAnalysis.Error = err.Error()
			fmt.Printf("  [Post %d/%d] ERROR: %v\n", i+1, len(posts), err)
		} else {
			posts[i].ValueAnalysis = *analysis
			fmt.Printf("  [Post %d/%d] OK - Tokens: %d - Time: %v\n",
				i+1, len(posts),
				analysis.Stats.TotalTokens,
				time.Since(postStart).Round(time.Millisecond))
		}

		cancel()
		time.Sleep(3 * time.Second)
	}

	if err := SavePostsToJson(filename, posts); err != nil {
		return fmt.Errorf("save json error: %w", err)
	}

	fmt.Printf("Completed in %v -> Saved to: %s\n", time.Since(startTime).Round(time.Millisecond), filename)
	return nil
}

func LoadStaticPosts(path string) ([]Post, error) {
	var posts []Post
	data, err := os.ReadFile(path)
	if err != nil {
		return posts, nil
	}

	err = json.Unmarshal(data, &posts)
	if err != nil {
		return posts, nil
	}
	return posts, nil
}

func LoadPostURLs(path string) (PostURLs, error) {
	var postURLs PostURLs
	data, err := os.ReadFile(path)
	if err != nil {
		return postURLs, nil
	}

	err = json.Unmarshal(data, &postURLs)
	if err != nil {
		return postURLs, nil
	}
	return postURLs, nil
}

func GetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("missing env: " + key)
	}
	return v
}

func GetOpenRouterClient() (AIClient, error) {
	key := os.Getenv("OPEN_ROUTER_KEY")
	if key == "" {
		return nil, fmt.Errorf("OPEN_ROUTER_KEY not set")
	}
	return NewOpenRouterClient(key), nil
}

func GetSiliconFlowClient() (AIClient, error) {
	apiKey := os.Getenv("SILICONFLOW_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SILICONFLOW_API_KEY not set")
	}
	return NewOpenAIClient(apiKey, "https://api.siliconflow.com/v1"), nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
