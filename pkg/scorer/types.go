package scorer

import "time"

var (
	GlobalStats Stats
)

type (
	Weights map[string]map[string]float64
	Models  map[string]string
)

type Stats struct {
	TotalCost       float64
	TotalTokens     int
	RequestCount    int
	TotalTime       time.Duration
	AvgCost         float64
	AvgTokens       float64
	AvgResponseTime time.Duration
}

type EmbedImage struct {
	URI string
	Alt string
}

type EmbedExternal struct {
	URI         string
	Title       string
	Description string
}

type FeedItem struct {
	URI           string         `json:"uri"`
	Text          string         `json:"text"`
	Replies       []string       `json:"replies"`
	Images        []EmbedImage   `json:"images"`
	External      *EmbedExternal `json:"external"`
	QuotedPostURI string         `json:"quoted_post_uri"`
	Values        SchwartzValues `json:"values"`
	ValuesArr     []int          `json:"values_arr"`
	Score         float64        `json:"score"`
	Model         string         `json:"model"`
	Stats         RequestStats   `json:"stats,omitempty"`
}

type RequestStats struct {
	ResponseTimeMs int64   `json:"response_time_ms"`
	TokensUsed     int     `json:"tokens_used"`
	CostUsd        float64 `json:"cost_usd"`
}

type Config struct {
	Weights Weights           `json:"weights"`
	Models  Models            `json:"models"`
	Ai      map[string]string `json:"ai"`
	Workers WorkerConfig      `json:"workers"`
}

type WorkerConfig struct {
	MaxConcurrent int `json:"max_concurrent"`
}

type SchwartzValues struct {
	SdThought     int `json:"sd_thought"`
	SdAction      int `json:"sd_action"`
	Stimulation   int `json:"stimulation"`
	Hedonism      int `json:"hedonism"`
	Achievement   int `json:"achievement"`
	Dominance     int `json:"dominance"`
	Resources     int `json:"resources"`
	Face          int `json:"face"`
	PersonalSec   int `json:"personal_sec"`
	SocietalSec   int `json:"societal_sec"`
	Tradition     int `json:"tradition"`
	RuleConf      int `json:"rule_conf"`
	InterConf     int `json:"inter_conf"`
	Humility      int `json:"humility"`
	Caring        int `json:"caring"`
	Dependability int `json:"dependability"`
	Universalism  int `json:"universalism"`
	Nature        int `json:"nature"`
	Tolerance     int `json:"tolerance"`
}
