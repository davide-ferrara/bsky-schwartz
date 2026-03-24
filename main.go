package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	openrouter "github.com/revrost/go-openrouter"
	"log"
	"os"
	"sort"
	"time"
)

type Stats struct {
	totalCost       float64
	totalTokens     int
	requestCount    int
	totalTime       time.Duration
	avgCost         float64
	avgTokens       float64
	avgResponseTime time.Duration
}

var stats Stats
var feed []FeedItem

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
	URI           string
	Text          string
	Replies       []string
	Images        []EmbedImage
	External      *EmbedExternal
	QuotedPostURI string
	HasVideo      bool
	IsRepost      bool
	RepostedBy    string
	Schwartz      []int
	Score         float64
	IndexedAt     time.Time
}

// calculateScore computes the final score based on Schwartz values and weights.
// Positive weights add directly: score += val * weight
// Negative weights are amplified by penalty: score += val * weight * penalty
//
// Example with penalty = 2.0:
//
//	val=6, weight=-5.0 → 6 * -5.0 * 2.0 = -60  (strong penalty)
//	val=3, weight=-5.0 → 3 * -5.0 * 2.0 = -30  (moderate penalty)
//	val=0, weight=-5.0 → 0 * -5.0 * 2.0 = 0    (neutral)
func (f *FeedItem) calculateScore(weights []float64) {
	var score float64
	penalty := 2.0

	for i, val := range f.Schwartz {
		weight := weights[i]
		if weight > 0 {
			score += float64(val) * weight
		} else if weight < 0 {
			score += float64(val) * weight * penalty
		}
	}

	f.Score = score
}

func hasURI(uri string) bool {
	for _, item := range feed {
		if item.URI == uri {
			return true
		}
	}
	return false
}

func addToFeed(items []FeedItem, model string, schwartzFile, basePrompt []byte, weights []float64) {
	for _, item := range items {
		if hasURI(item.URI) {
			fmt.Printf("    [SKIP] Duplicate: %s\n", item.URI)
			continue
		}

		prompt := fmt.Sprintf("%s\n\n---BEGIN POST---\n%s\n---END POST---\n%s", schwartzFile, item.Text, basePrompt)

		var err error
		item.Schwartz, err = calculateValues(model, prompt)
		if err != nil {
			fmt.Printf("    [ERROR] %s: %v\n", item.URI, err)
			continue
		}

		item.calculateScore(weights)
		feed = append(feed, item)
		fmt.Printf("    [+] Score: %.2f | %s\n", item.Score, truncate(item.Text, 60))
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func printFeed() {
	sort.Slice(feed, func(i, j int) bool {
		return feed[i].Score > feed[j].Score
	})

	for i, item := range feed {
		fmt.Printf("[%d] Score: %.2f | %s\n", i+1, item.Score, item.Text)
	}
}

func saveFeed(filename string, model string) error {
	sort.Slice(feed, func(i, j int) bool {
		return feed[i].Score > feed[j].Score
	})

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "Model: %s\n", model)
	fmt.Fprintf(f, "Generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f, "Posts: %d\n", len(feed))
	fmt.Fprintf(f, "--- Stats ---\n")
	fmt.Fprintf(f, "Requests: %d\n", stats.requestCount)
	fmt.Fprintf(f, "Total tokens: %d\n", stats.totalTokens)
	fmt.Fprintf(f, "Total cost: $%.6f\n", stats.totalCost)
	fmt.Fprintf(f, "Avg tokens per request: %.2f\n", float64(stats.totalTokens)/float64(stats.requestCount))
	fmt.Fprintf(f, "Avg cost per request: $%.6f\n", stats.totalCost/float64(stats.requestCount))
	fmt.Fprintf(f, "Total time: %v\n", stats.totalTime)
	fmt.Fprintf(f, "Avg response time: %v\n", stats.totalTime/time.Duration(stats.requestCount))
	fmt.Fprintf(f, "\n")

	for i, item := range feed {
		fmt.Fprintf(f, "=== POST #%d ===\n", i+1)
		fmt.Fprintf(f, "Score: %.2f\n", item.Score)
		fmt.Fprintf(f, "URI: %s\n", item.URI)
		fmt.Fprintf(f, "Text: %s\n", item.Text)
		fmt.Fprintf(f, "IndexedAt: %s\n", item.IndexedAt)

		if len(item.Schwartz) > 0 {
			fmt.Fprintf(f, "Schwartz Values:\n")
			for j, val := range item.Schwartz {
				fmt.Fprintf(f, "  %2d. %-25s %d\n", j, schwartzNames[j], val)
			}
			fmt.Fprintf(f, "  Array: %v\n", item.Schwartz)
		}

		if len(item.Replies) > 0 {
			fmt.Fprintf(f, "Replies:\n")
			for _, r := range item.Replies {
				fmt.Fprintf(f, "  - %s\n", r)
			}
		}

		if len(item.Images) > 0 {
			fmt.Fprintf(f, "Images:\n")
			for _, img := range item.Images {
				fmt.Fprintf(f, "  - %s (alt: %s)\n", img.URI, img.Alt)
			}
		}

		if item.External != nil {
			fmt.Fprintf(f, "External: %s - %s\n", item.External.Title, item.External.URI)
		}

		if item.QuotedPostURI != "" {
			fmt.Fprintf(f, "QuotedPost: %s\n", item.QuotedPostURI)
		}

		if item.HasVideo {
			fmt.Fprintf(f, "HasVideo: true\n")
		}

		if item.IsRepost {
			fmt.Fprintf(f, "RepostedBy: %s\n", item.RepostedBy)
		}

		fmt.Fprintf(f, "\n")
	}

	return nil
}

var keyToIndex = map[string]int{
	"sd_thought":    0,
	"sd_action":     1,
	"stimulation":   2,
	"hedonism":      3,
	"achievement":   4,
	"dominance":     5,
	"resources":     6,
	"face":          7,
	"personal_sec":  8,
	"societal_sec":  9,
	"tradition":     10,
	"rule_conf":     11,
	"inter_conf":    12,
	"humility":      13,
	"caring":        14,
	"dependability": 15,
	"universalism":  16,
	"nature":        17,
	"tolerance":     18,
}

var schwartzNames = []string{
	"Self-directed thoughts",
	"Self-directed actions",
	"Stimulation",
	"Hedonism",
	"Achievement",
	"Dominance",
	"Resources",
	"Face",
	"Personal security",
	"Societal security",
	"Tradition",
	"Rule conformity",
	"Interpersonal conformity",
	"Humility",
	"Caring",
	"Dependability",
	"Universal concern",
	"Preservation of nature",
	"Tolerance",
}

func getDefaultWeights() []float64 {
	return []float64{
		// OPENNESS TO CHANGE (Neutrale/Positivo)
		0,    // 0. Self-directed thoughts
		0,    // 1. Self-directed actions
		-1.0, // 2. Stimulation
		-5.0, // 3. Hedonism (Spesso legato a retorica di consumo)

		// SELF-ENHANCEMENT (Forte Penalizzazione)
		-5.0, // 4. Achievement (Successo sociale egoistico)
		-5.0, // 5. Dominance (Il valore chiave dell'odio/potere)
		-5.0, // 6. Resources (Controllo materiale)
		-4.0, // 7. Face (Mantenimento immagine pubblica/onore)

		// CONSERVATION (Forte Penalizzazione per la retorica identitaria)
		-1.0, // 8. Personal Security
		2.0,  // 9. Societal Security (Sicurezza usata come "muro")
		-5.0, // 10. Tradition (Usata come esclusione)
		0.0,  // 11. Rule Conformity (Obbedienza cieca)
		5.0,  // 12. Interpersonal Conformity
		5.0,  // 13. Humility

		// SELF-TRANSCENDENCE (Massimo Premio)
		5.0, // 14. Caring
		5.0, // 15. Dependability
		5.0, // 16. Universal Concern (Uguaglianza/Giustizia)
		5.0, // 17. Preservation of Nature
		5.0, // 18. Tolerance (Accettazione del diverso)
	}
}

func calculateValues(model string, prompt string) ([]int, error) {
	client := openrouter.NewClient(
		os.Getenv("OPEN_ROUTER_KEY"),
	)

	start := time.Now()
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Model:       model,
			Temperature: 0,
			Messages: []openrouter.ChatCompletionMessage{
				openrouter.UserMessage(prompt),
			},
		},
	)
	elapsed := time.Since(start)

	if err != nil {
		return []int{}, fmt.Errorf("ChatCompletion error: %v", err)
	}

	stats.totalCost += resp.Usage.Cost
	stats.totalTokens += resp.Usage.PromptTokens + resp.Usage.CompletionTokens
	stats.totalTime += elapsed
	stats.requestCount++

	var jsonMap map[string]int
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content.Text), &jsonMap)
	if err != nil {
		return []int{}, fmt.Errorf("JSON parse error: %v", err)
	}

	arr := make([]int, 19)
	for key, val := range jsonMap {
		if idx, ok := keyToIndex[key]; ok {
			arr[idx] = val
		}
	}

	return arr, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env key, exiting...")
		os.Exit(0)
	}

	model := "qwen/qwen-2.5-72b-instruct"
	// model := "openai/gpt-4o-mini"

	schwartzFile, err := os.ReadFile("SCHWARTZ.md")
	if err != nil {
		log.Fatal("Could not load SCHWARTZ.md")
	}

	basePrompt, err := os.ReadFile("PROMPT.md")
	if err != nil {
		log.Fatal("Could not load PROMPT.md")
	}

	weights := getDefaultWeights()

	queries := []string{
		"Gaza",
		"Genocidio",
		"7 Ottobre",
		"Netanyahu",
	}

	fmt.Printf("=== Schwartz Feed Generator ===\n")
	fmt.Printf("Model: %s\n", model)
	fmt.Printf("Queries: %d\n\n", len(queries))

	client := bClient{}.new(os.Getenv("BLUESKY_HANDLE"), os.Getenv("BLUESKY_APP_PASSWORD"))

	for i, q := range queries {
		fmt.Printf("[%d/%d] Query: \"%s\"\n", i+1, len(queries), q)
		items := client.queryPost(q, 5)
		fmt.Printf("  Found: %d posts\n", len(items))
		addToFeed(items, model, schwartzFile, basePrompt, weights)
	}

	fmt.Printf("\n=== Feed Summary ===\n")
	fmt.Printf("Total posts analyzed: %d\n", len(feed))
	fmt.Printf("AI requests: %d\n", stats.requestCount)
	fmt.Printf("Total cost: $%.6f\n", stats.totalCost)
	fmt.Printf("Total time: %v\n\n", stats.totalTime)

	printFeed()

	err = saveFeed(fmt.Sprintf("feed_%s.txt", time.Now().Format("20060102_150405")), model)
	if err != nil {
		panic(0)
	}

	fmt.Printf("\nFeed saved to file.\n")
}
