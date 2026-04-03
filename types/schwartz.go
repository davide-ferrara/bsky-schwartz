package types

type SchwartzValues map[string]int

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
