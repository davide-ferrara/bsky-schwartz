package scorer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	openrouter "github.com/revrost/go-openrouter"
)

var (
	cfg  *Config
	once sync.Once
)

func Init() error {
	var err error
	once.Do(func() {
		cfg, err = loadConfig()
		if err != nil {
			return
		}
		initCachedReader()
	})
	return err
}

func GetConfig() *Config {
	return cfg
}

func loadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.json"
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

func getORClient() *openrouter.Client {
	key := os.Getenv("OPEN_ROUTER_KEY")
	if key == "" {
		panic("OPEN_ROUTER_KEY is empty")
	}
	return openrouter.NewClient(key)
}

func stripMarkdownCodeFences(data []byte) []byte {
	text := string(data)
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSuffix(text, "```json")
	return []byte(strings.TrimSpace(text))
}

func (f *FeedItem) generatePrompt() (string, error) {
	config := GetConfig()

	promptData, err := GetCachedReader()(config.Ai["prompt"])
	if err != nil {
		return "", fmt.Errorf("error reading PROMPT.MD: %w", err)
	}

	schwartzData, err := GetCachedReader()(config.Ai["schwartz"])
	if err != nil {
		return "", fmt.Errorf("error reading SCHWARTZ.MD: %w", err)
	}

	return fmt.Sprintf("%s\n---BEGIN POST---\n%s\n---END POST---\n%s", schwartzData, f.Text, promptData), nil
}

func (f *FeedItem) ValueAlignment(model string) error {
	log := slog.With("post_uri", f.URI, "model", model)
	log.Info("ai analysis started")

	start := time.Now()

	prompt, err := f.generatePrompt()
	if err != nil {
		log.Error("failed to generate prompt", "error", err)
		return err
	}

	client := getORClient()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	aiResp, err := client.CreateChatCompletion(
		ctx,
		openrouter.ChatCompletionRequest{
			Model:       model,
			Temperature: 0,
			Messages: []openrouter.ChatCompletionMessage{
				openrouter.UserMessage(prompt),
			},
		},
	)
	if err != nil {
		log.Error("openrouter client error", "error", err)
		return fmt.Errorf("openrouter client error: %v", err)
	}

	elapsed := time.Since(start)

	data := []byte(aiResp.Choices[0].Message.Content.Text)

	data = stripMarkdownCodeFences(data)

	err = json.Unmarshal(data, &f.Values)
	if err != nil {
		log.Error("failed to unmarshal ai response", "error", err, "response", string(data))
		return fmt.Errorf("unmarshal error: %v, data: %s", err, string(data))
	}

	f.ValuesArr = f.Values.ToArray()
	f.calculateScore()
	f.Model = model

	f.Stats = RequestStats{
		ResponseTimeMs: elapsed.Milliseconds(),
		TokensUsed:     aiResp.Usage.TotalTokens,
		CostUsd:        aiResp.Usage.Cost,
	}

	log.Info("ai analysis completed",
		"response_time_ms", elapsed.Milliseconds(),
		"tokens_used", aiResp.Usage.TotalTokens,
		"cost_usd", aiResp.Usage.Cost,
		"score", f.Score,
	)

	return nil
}

func (f *FeedItem) calculateScore() {
	weights := weightsToArray()
	values := f.ValuesArr

	var score float64
	for i := range values {
		score += float64(values[i]) * weights[i]
	}

	f.Score = score
}

// func (f *FeedItem) CalculateScoreCustom(weights []float64) {
// 	var score float64
// 	penalty := 2.0
//
// 	for i, val := range f.ValuesArr {
// 		weight := weights[i]
// 		if weight > 0 {
// 			score += float64(val) * weight
// 		} else if weight < 0 {
// 			score += float64(val) * weight * penalty
// 		}
// 	}
//
// 	f.Score = score
// }

func parseSchwartzValues(data []byte) (*SchwartzValues, error) {
	var v SchwartzValues
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (v *SchwartzValues) ToArray() []int {
	return []int{
		v.SdThought, v.SdAction, v.Stimulation, v.Hedonism,
		v.Achievement, v.Dominance, v.Resources, v.Face,
		v.PersonalSec, v.SocietalSec, v.Tradition, v.RuleConf,
		v.InterConf, v.Humility, v.Caring, v.Dependability,
		v.Universalism, v.Nature, v.Tolerance,
	}
}

func weightsToArray() []float64 {
	cfg := GetConfig()
	w := cfg.Weights["left"]
	return []float64{
		w["Self-directed thoughts"],
		w["Self-directed actions"],
		w["Stimulation"],
		w["Hedonism"],
		w["Achievement"],
		w["Dominance"],
		w["Resources"],
		w["Face"],
		w["Personal security"],
		w["Societal security"],
		w["Tradition"],
		w["Rule conformity"],
		w["Interpersonal conformity"],
		w["Humility"],
		w["Caring"],
		w["Dependability"],
		w["Universal concern"],
		w["Preservation of nature"],
		w["Tolerance"],
	}
}
