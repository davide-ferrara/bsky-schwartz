package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/joho/godotenv"
	openrouter "github.com/revrost/go-openrouter"
)

type Record struct {
	Type  string   `json:"type"`
	Text  string   `json:"text"`
	Langs []string `json:"langs"`
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
	// Jetstream() // JetStream firehose
	// Indigo()    // Search posts
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env key, exiting...")
		os.Exit(0)
	}

	// token := countToken("SCHWARTZ.md")
	// log.Printf("SCHWARTZ.md Tokens are approx: %v\n", token)
	// log.Printf("API Call cost is: %.6f$\n", AiCost(token, 0.12))

	schwartzDef, err := os.ReadFile("SCHWARTZ.md")
	if err != nil {
		log.Fatal("Could not load Schwartz.md")
	}

	basePrompt := "Based on the Schwartz Definitions, judge the following Social Media Post. IMPORTANT: Respond with ONLY an array of 19 integers (scores 0-6). Example: [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19]. No explanation, no JSON, no text before or after. Start with [ and end with ]."

	for i := 1; i <= 5; i++ {
		filename := fmt.Sprintf("post_collections/post_%v.txt", i)
		post, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("Skipping: %v\n", filename)
			continue
		}

		prompt := fmt.Sprintf("%s\n\n---BEGIN POST---\n%s\n---END POST---\n\nSchwartz Definitions:\n%s", basePrompt, post, schwartzDef)

		values, err := ai(prompt)
		if err != nil {
			log.Printf("Error processing %s: %v\n", filename, err)
			continue
		}

		// // For Debugging
		// sz := SchwartzValues{}
		// sz.OpennessToChange.SelfDirectedThoughts = values[0]
		// sz.OpennessToChange.SelfDirectedActions = values[1]
		// sz.OpennessToChange.Stimulation = values[2]
		// sz.OpennessToChange.Hedonism = values[3]
		// sz.SelfEnhancement.Achievement = values[4]
		// sz.SelfEnhancement.Dominance = values[5]
		// sz.SelfEnhancement.Resources = values[6]
		// sz.SelfEnhancement.Face = values[7]
		// sz.Conservation.PersonalSecurity = values[8]
		// sz.Conservation.SocietalSecurity = values[9]
		// sz.Conservation.Tradition = values[10]
		// sz.Conservation.RuleConformity = values[11]
		// sz.Conservation.InterpersonalConformity = values[12]
		// sz.Conservation.Humility = values[13]
		// sz.SelfTranscendence.Caring = values[14]
		// sz.SelfTranscendence.Dependability = values[15]
		// sz.SelfTranscendence.UniversalConcern = values[16]
		// sz.SelfTranscendence.PreservationOfNature = values[17]
		// sz.SelfTranscendence.Tolerance = values[18]
		//
		// szJSON, _ := json.MarshalIndent(sz, "", "  ")
		// fmt.Printf("[%s]\n%s\n", filename, string(szJSON))

		// Extreme Left weights [-1, 1] step 0.25
		weights := []float64{
			0.5,   // 0.  Self-directed thoughts
			0.5,   // 1.  Self-directed actions
			0.0,   // 2.  Stimulation
			-0.25, // 3.  Hedonism
			-0.5,  // 4.  Achievement
			-1.0,  // 5.  Dominance
			-0.5,  // 6.  Resources
			-1.0,  // 7.  Face
			-0.25, // 8.  Personal Security
			0.0,   // 9.  Societal Security
			-1.0,  // 10. Tradition
			-0.5,  // 11. Rule Conformity
			1.0,   // 12. Interpersonal Conformity
			1.0,   // 13. Humility
			1.0,   // 14. Caring
			1.0,   // 15. Dependability
			1.0,   // 16. Universal Concern
			1.0,   // 17. Preservation of Nature
			1.0,   // 18. Tolerance
		}

		var score float64
		for j := 0; j < 19; j++ {
			score += float64(values[j]) * weights[j]
		}
		normalized := score / 114.0

		fmt.Printf("Score: %.2f (normalized: %.2f)\n\n", score, normalized)
	}
}

func AiCost(tokens float64, price float64) float64 {
	return tokens * (price / 1000000)
}

func countToken(filename string) float64 {
	fp, err := os.OpenFile(filename, 0, 0644)
	if err != nil {
		log.Fatal("Could not open file!")
	}

	stat, err := fp.Stat()
	if err != nil {
		log.Fatal("Could not stat!")
	}

	return math.Round(float64(stat.Size()) / 4)
}

func ai(prompt string) ([]int, error) {
	client := openrouter.NewClient(
		os.Getenv("OPEN_ROUTER_KEY"),
	)

	// model := "qwen/qwen-2.5-72b-instruct"
	model := "openai/gpt-4o-mini"

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Model: model,
			Messages: []openrouter.ChatCompletionMessage{
				openrouter.UserMessage(prompt),
			},
		},
	)

	if err != nil {
		return []int{}, fmt.Errorf("ChatCompletion error: %v", err)
	}

	var arr []int
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content.Text), &arr)
	if err != nil {
		return []int{}, fmt.Errorf("JSON parse error: %v", err)
	}
	return arr, nil
}
