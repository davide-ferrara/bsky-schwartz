package main

import (
	"context"

	"github.com/bluesky-social/indigo/xrpc"
	"github.com/openai/openai-go/v3"
	openrouter "github.com/revrost/go-openrouter"
)

type Client struct {
	client *xrpc.Client
}

type AIClient interface {
	CreateChatCompletion(ctx context.Context, model string, prompt string) (*AIResponse, error)
}

type AIResponse struct {
	Content          string
	Model            string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	CostUsd          float64
	Provider         string
}

type OpenRouterClient struct {
	client *openrouter.Client
}

type OpenAIClient struct {
	client *openai.Client
}
