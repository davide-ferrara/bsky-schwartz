package scorer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPostScoreHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/post/score/:model", PostScoreHanlder)

	tests := []struct {
		name       string
		model      string
		wantStatus int
	}{
		{"valid_gpt", "gpt", http.StatusOK},
		{"valid_qwen", "qwen", http.StatusOK},
		{"invalid_model", "invalid", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/post/score/"+tt.model, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestParseSchwartzValues(t *testing.T) {
	jsonData := `{
		"sd_thought": 5,
		"sd_action": 3,
		"stimulation": 7,
		"hedonism": 2,
		"achievement": 8,
		"dominance": 4,
		"resources": 6,
		"face": 1,
		"personal_sec": 9,
		"societal_sec": 3,
		"tradition": 2,
		"rule_conf": 4,
		"inter_conf": 5,
		"humility": 6,
		"caring": 7,
		"dependability": 8,
		"universalism": 9,
		"nature": 4,
		"tolerance": 5
	}`

	v, err := parseSchwartzValues([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseSchwartzValues failed: %v", err)
	}

	if v.SdThought != 5 {
		t.Errorf("SdThought = %d, want 5", v.SdThought)
	}
	if v.SdAction != 3 {
		t.Errorf("SdAction = %d, want 3", v.SdAction)
	}
	if v.Stimulation != 7 {
		t.Errorf("Stimulation = %d, want 7", v.Stimulation)
	}

	arr := v.ToArray()
	if len(arr) != 19 {
		t.Errorf("ToArray length = %d, want 19", len(arr))
	}
	if arr[0] != 5 || arr[1] != 3 || arr[2] != 7 {
		t.Errorf("ToArray = %v, want first 3 elements to be [5, 3, 7]", arr[:3])
	}

	fmt.Printf("Parsed SchwartzValues: %+v\n", v)
	fmt.Printf("Array: %v\n", arr)
}
