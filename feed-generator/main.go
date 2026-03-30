package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/revrost/go-openrouter"
)

func main() {
	flag.Parse()

	if err := godotenv.Load(".env"); err != nil {
		fmt.Fprintf(os.Stderr, "warning: .env not found: %v\n", err)
	}

	openrouterClient, err := GetAiClient()
	if err != nil {
		panic("failed to init AI Client: " + err.Error())
	}

	models := []string{"openai/gpt-4.1-mini", "google/gemini-3.1-flash-lite-preview", "mistralai/mistral-small-2603", "qwen/qwen3.5-9b"}

	// posts := DownloadFeed("./feed/conservation.json")

	posts, err := LoadStaticPosts("./feed_20260330212226.json")
	if err != nil {
		fmt.Println("Could not open JSON")
		panic(err)
	}

	fmt.Println("Starting analysis...")

	limit := 10

	if err := runAnalysisSync(openrouterClient, posts[:limit], models[3]); err != nil {
		fmt.Println("Error running analysis with", models[3], ":", err)
	}
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

func runAnalysisSync(c *openrouter.Client, posts []Post, model string) error {
	fmt.Println("Started post analysis...")

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	for i := range posts {
		fmt.Printf("Calculating rating with model %v for post: %s\n", model, posts[i].AtURI)

		analysis, err := CalculateRating(ctx, c, model, &posts[i])
		if err != nil {
			posts[i].ValueAnalysis.Error = err.Error()
			fmt.Println("Error in ValueAnalysis:", err)
		} else {
			posts[i].ValueAnalysis = *analysis
		}

		time.Sleep(5 * time.Second)
	}

	re := regexp.MustCompile(`^[^/]+/(.+)$`)
	matches := re.FindStringSubmatch(model)
	model = matches[1]

	filename := fmt.Sprintf("post_%s", model)
	if err := SavePostsToJson(filename, posts); err != nil {
		return fmt.Errorf("save json error: %w", err)
	}

	fmt.Println("Analysis completed and saved to:", filename)
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

func GetAiClient() (*openrouter.Client, error) {
	key := os.Getenv("OPEN_ROUTER_KEY")
	if key == "" {
		return &openrouter.Client{}, fmt.Errorf("OPEN_ROUTER_KEY not set")
	}
	return openrouter.NewClient(key), nil
}
