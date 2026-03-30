package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/revrost/go-openrouter"
)

func BuildPromptContent(post *Post) string {
	content := map[string]interface{}{
		"text": post.Text,
	}

	if len(post.Langs) > 0 {
		content["language"] = strings.Join(post.Langs, ", ")
	}

	if len(post.Links) > 0 {
		var links []map[string]string
		for _, l := range post.Links {
			links = append(links, map[string]string{
				"title":       l.Title,
				"description": l.Description,
			})
		}
		content["external_links"] = links
	}

	jsonBytes, _ := json.Marshal(content)
	return fmt.Sprintf("<post>\n%s\n</post>", jsonBytes)
}

func CalculateRating(ctx context.Context, client *openrouter.Client, model string, post *Post) (*ValueAnalysis, error) {
	start := time.Now()

	taskPrompt, err := os.ReadFile("./prompts/PROMPT_V3.md")
	if err != nil {
		return nil, fmt.Errorf("reading prompt: %w", err)
	}

	promptContent := BuildPromptContent(post)
	prompt := fmt.Sprintf("%s\n\n%s", string(taskPrompt), promptContent)

	resp, err := client.CreateChatCompletion(
		ctx,
		openrouter.ChatCompletionRequest{
			Model:       model,
			Temperature: 0.3,
			Messages: []openrouter.ChatCompletionMessage{
				openrouter.UserMessage(prompt),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("openrouter error: %w", err)
	}

	elapsed := time.Since(start)

	jsonStr := resp.Choices[0].Message.Content.Text

	var result struct {
		Rating    SchwartzValues `json:"Rating"`
		Reasoning string         `json:"Reasoning"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w, json: %s", err, jsonStr)
	}

	stats := AIStats{
		Model:            resp.Model,
		ResponseTimeMs:   elapsed.Milliseconds(),
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
		TotalTokens:      resp.Usage.TotalTokens,
		CostUsd:          resp.Usage.Cost,
		Provider:         resp.Provider,
	}

	return &ValueAnalysis{
		Rating:    result.Rating,
		Reasoning: result.Reasoning,
		Stats:     stats,
	}, nil
}
