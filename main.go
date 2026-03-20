package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	openrouter "github.com/revrost/go-openrouter"
	"log"
	"os"
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

type Record struct {
	Type  string   `json:"type"`
	Text  string   `json:"text"`
	Langs []string `json:"langs"`
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

type SchwartzValues struct {
	OpennessToChange struct {
		SelfDirectedThoughts int `json:"self_directed_thoughts"`
		SelfDirectedActions  int `json:"self_directed_actions"`
		Stimulation          int `json:"stimulation"`
		Hedonism             int `json:"hedonism"`
	} `json:"openness_to_change"`
	SelfEnhancement struct {
		Achievement int `json:"achievement"`
		Dominance   int `json:"dominance"`
		Resources   int `json:"resources"`
		Face        int `json:"face"`
	} `json:"self_enhancement"`
	Conservation struct {
		PersonalSecurity        int `json:"personal_security"`
		SocietalSecurity        int `json:"societal_security"`
		Tradition               int `json:"tradition"`
		RuleConformity          int `json:"rule_conformity"`
		InterpersonalConformity int `json:"interpersonal_conformity"`
		Humility                int `json:"humility"`
	} `json:"conservation"`
	SelfTranscendence struct {
		Caring               int `json:"caring"`
		Dependability        int `json:"dependability"`
		UniversalConcern     int `json:"universal_concern"`
		PreservationOfNature int `json:"preservation_of_nature"`
		Tolerance            int `json:"tolerance"`
	} `json:"self_transcendence"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env key, exiting...")
		os.Exit(0)
	}

	// model := "qwen/qwen-2.5-72b-instruct"
	// // model := "deepseek/deepseek-v3.2"
	// // model := "openai/gpt-4o-mini"
	//
	//
	// schwartzFile, err := os.ReadFile("SCHWARTZ.md")
	// if err != nil {
	// 	log.Fatal("Could not load SCHWARTZ.md")
	// }
	//
	// basePrompt, err := os.ReadFile("PROMPT.md")
	// if err != nil {
	// 	log.Fatal("Could not load PROMPT.md")
	// }
	//
	// weights := []float64{
	// 	// OPENNESS TO CHANGE (Neutrale/Positivo)
	// 	0,    // 0. Self-directed thoughts
	// 	0,    // 1. Self-directed actions
	// 	-1.0, // 2. Stimulation
	// 	-5.0, // 3. Hedonism (Spesso legato a retorica di consumo)
	//
	// 	// SELF-ENHANCEMENT (Forte Penalizzazione)
	// 	-5.0, // 4. Achievement (Successo sociale egoistico)
	// 	-5.0, // 5. Dominance (Il valore chiave dell'odio/potere)
	// 	-5.0, // 6. Resources (Controllo materiale)
	// 	-4.0, // 7. Face (Mantenimento immagine pubblica/onore)
	//
	// 	// CONSERVATION (Forte Penalizzazione per la retorica identitaria)
	// 	-1.0, // 8. Personal Security
	// 	-4.0, // 9. Societal Security (Sicurezza usata come "muro")
	// 	-5.0, // 10. Tradition (Usata come esclusione)
	// 	-3.0, // 11. Rule Conformity (Obbedienza cieca)
	// 	5.0,  // 12. Interpersonal Conformity
	// 	5.0,  // 13. Humility
	//
	// 	// SELF-TRANSCENDENCE (Massimo Premio)
	// 	5.0, // 14. Caring
	// 	5.0, // 15. Dependability
	// 	5.0, // 16. Universal Concern (Uguaglianza/Giustizia)
	// 	5.0, // 17. Preservation of Nature
	// 	5.0, // 18. Tolerance (Accettazione del diverso)
	// }

	// runTest(model, schwartzFile, basePrompt, weights)

	query := "Giorgia Meloni Referendum"
	queryPost(os.Getenv("BLUESKY_HANDLE"), os.Getenv("BLUESKY_APP_PASSWORD"), query, 3)
}

func queryAI(model string, prompt string) ([]int, error) {
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

	fmt.Printf("AI Response: %s\n\n", resp.Choices[0].Message.Content.Text)

	stats.totalCost += resp.Usage.Cost
	stats.totalTokens += resp.Usage.PromptTokens + resp.Usage.CompletionTokens
	stats.totalTime += elapsed
	stats.requestCount++

	fmt.Printf("Input tokens: %d | Output tokens: %d | Cost: $%.6f | Response time: %v | Total: $%.6f\n",
		resp.Usage.PromptTokens,
		resp.Usage.CompletionTokens,
		resp.Usage.Cost,
		elapsed,
		stats.totalCost,
	)

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

func runTest(model string, schwartzFile, basePrompt []byte, weights []float64) {
	numOfPosts := 5
	for i := 1; i <= numOfPosts; i++ {
		filename := fmt.Sprintf("post_collections/post_%v.txt", i)
		post, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("Skipping: %v\n", filename)
			continue
		}

		prompt := fmt.Sprintf("%s\n\n---BEGIN POST---\n%s\n---END POST---\n%s", schwartzFile, post, basePrompt)

		values, err := queryAI(model, prompt)
		if err != nil {
			log.Printf("Error processing %s: %v\n", filename, err)
			continue
		}

		// TODO: Separate into funv
		var score float64
		for i := 0; i < 19; i++ {
			val := float64(values[i])
			weight := weights[i]
			penalty := 3.0

			if weight > 0 {
				// PREMIO PURO: Se il valore c'è, aggiungi punti. Se è 0, non succede nulla.
				// Questo evita che un post di sinistra "cada" perché non parla di natura.
				score += val * weight
			} else if weight < 0 {
				// PENALITÀ ATTIVA: Usiamo la traslazione solo qui.
				// Se il valore odiato è 0, (0-3) * -5 = +15 (un premio per non essere d'odio).
				// Se il valore odiato è 6, (6-3) * -5 = -15 (punizione).
				score += (val - penalty) * weight
			}
		}

		fmt.Printf("Score: %.2f \n\n", score)
	}

	stats.avgCost = stats.totalCost / float64(stats.requestCount)
	stats.avgTokens = float64(stats.totalTokens) / float64(stats.requestCount)
	stats.avgResponseTime = stats.totalTime / time.Duration(stats.requestCount)

	weightsID := hashWeights(weights)

	fmt.Printf("=== Session Stats ===\n")
	fmt.Printf("Model: %s\n", model)
	fmt.Printf("Weights ID: %s\n", weightsID)
	fmt.Printf("Requests: %d\n", stats.requestCount)
	fmt.Printf("Total cost: $%.6f\n", stats.totalCost)
	fmt.Printf("Avg cost per request: $%.6f\n", stats.avgCost)
	fmt.Printf("Total tokens: %d\n", stats.totalTokens)
	fmt.Printf("Avg tokens per request: %.2f\n", stats.avgTokens)
	fmt.Printf("Total time: %v\n", stats.totalTime)
	fmt.Printf("Avg response time: %v\n", stats.avgResponseTime)
}

// TODO: fix scoreRange()
func scoreRange(weights []float64) (min, max float64) {
	for _, w := range weights {
		if w > 0 {
			max += w * 6
		} else if w < 0 {
			min += w * 6
		}
	}
	return min, max
}

func hashWeights(weights []float64) string {
	data, _ := json.Marshal(weights)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
