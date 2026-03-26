package main

import (
	"encoding/json"
	"fmt"
	"os"

	"bsky-schwarz/pkg/bluesky"
	"bsky-schwarz/pkg/scorer"

	"github.com/joho/godotenv"
)

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

	client := bluesky.NewClient(handle, appPassword)

	feed := client.QueryPosts("Gaza", 1)

	model := scorer.GetConfig().Models["gpt"]
	err = feed[0].ValueAlignment(model)
	if err != nil {
		panic(err)
	}

	fmt.Println(feed[0].Text)
	b, _ := json.MarshalIndent(feed[0].Values, "", "  ")
	fmt.Println(string(b))
}

func getEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	panic("missing environment variable: " + key)
}
