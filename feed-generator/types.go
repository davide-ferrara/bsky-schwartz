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

type PostImage struct {
	Alt   string
	Image string
}

type PostLink struct {
	Uri         string
	Title       string
	Description string
	Thumb       string
}

type PostFacet struct {
	Type  string
	Value string
}

type Post struct {
	URL           string
	AtURI         string
	Text          string
	CreatedAt     string
	Labels        []string
	Langs         []string
	Tags          []string
	Images        []PostImage
	Links         []PostLink
	Facets        []PostFacet
	AuthorName    string
	ReplyRoot     string
	ReplyParent   string
	ValueAnalysis ValueAnalysis
}

type SchwartzValues struct {
	Reputation          int `json:"Reputation"`
	Power               int `json:"Power"`
	Wealth              int `json:"Wealth"`
	Achievement         int `json:"Achievement"`
	Pleasure            int `json:"Pleasure"`
	IndependentThoughts int `json:"Independent thoughts"`
	IndependentActions  int `json:"Independent actions"`
	Stimulation         int `json:"Stimulation"`
	PersonalSecurity    int `json:"Personal security"`
	SocietalSecurity    int `json:"Societal security"`
	Tradition           int `json:"Tradition"`
	Lawfulness          int `json:"Lawfulness"`
	Respect             int `json:"Respect"`
	Humility            int `json:"Humility"`
	Responsibility      int `json:"Responsibility"`
	Caring              int `json:"Caring"`
	Equality            int `json:"Equality"`
	Nature              int `json:"Nature"`
	Tolerance           int `json:"Tolerance"`
}

type AIStats struct {
	Model            string  `json:"model"`
	ResponseTimeMs   int64   `json:"response_time_ms"`
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	CostUsd          float64 `json:"cost_usd"`
	Provider         string  `json:"provider"`
}

type ValueAnalysis struct {
	Rating    SchwartzValues `json:"Rating"`
	Reasoning string         `json:"Reasoning"`
	Score     int            `json:"Score"`
	Stats     AIStats        `json:"Stats"`
	Error     string         `json:"error,omitempty"`
}

type PostURLs struct {
	Urls []string `json:"urls"`
}
