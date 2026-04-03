package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"bsky-schwartz/internal/models"
	"bsky-schwartz/types"

	"github.com/bluesky-social/indigo/api/bsky"
)

// AlgoHandler - Firma funzione per un algoritmo di feed
type AlgoHandler func(ctx context.Context, limit int, cursor string, userDID string) (*bsky.FeedGetFeedSkeleton_Output, error)

// Algos - Registro degli algoritmi disponibili
var Algos = map[string]AlgoHandler{
	"values": feedValueBased,
}

// =============================================================================
// POST DATABASE
// =============================================================================

var (
	postsOnce sync.Once
	posts     []types.Post
	postsErr  error
)

// LoadFeed - Carica i post da feed_test.json (singleton)
func LoadFeed() ([]types.Post, error) {
	postsOnce.Do(func() {
		data, err := os.ReadFile("./feed_test.json")
		if err != nil {
			postsErr = err
			return
		}
		if err := json.Unmarshal(data, &posts); err != nil {
			postsErr = err
			return
		}
	})
	return posts, postsErr
}

// =============================================================================
// SCORING
// =============================================================================

func jsonKeyToWeightKey(jsonKey string) string {
	return strings.ToLower(strings.ReplaceAll(jsonKey, " ", "_"))
}

// CalculateScore - Calcola lo score di un post rispetto ai weights
func CalculateScore(posts []types.Post, weights map[string]float64) {
	for i := range posts {
		var score int
		for key, rating := range posts[i].ValueAnalysis.Rating {
			score += rating * int(weights[strings.ToLower(key)])
		}
		posts[i].ValueAnalysis.Score = score
		fmt.Println("DEBUG: ", posts[i].ValueAnalysis.Score)
	}
}

// =============================================================================
// FEED GENERATOR
// =============================================================================

// feedValueBased - Restituisce i post ordinati per score
func feedValueBased(ctx context.Context, limit int, cursor string, userDID string) (*bsky.FeedGetFeedSkeleton_Output, error) {
	// Recupera weights utente dalla memoria condivisa
	weights := GlobalUserWeights.Get(userDID)
	if weights == nil {
		// Weight neutri (tutti a 0) - feed non personalizzato
		weights = make(map[string]float64)
		for _, v := range models.SwartzValues {
			weights[v.ID] = 0.0
		}
		fmt.Printf("ℹ️ Nessun weights per DID=%s, uso pesi neutri\n", userDID)
	} else {
		fmt.Printf("✓ Weights trovati per DID=%s (%d valori)\n", userDID, len(weights))
	}

	// Carica i post dal file (singleton)
	posts, err := LoadFeed()
	if err != nil {
		return nil, fmt.Errorf("loading feed: %w", err)
	}

	CalculateScore(posts, weights)

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ValueAnalysis.Score > posts[j].ValueAnalysis.Score
	})

	// Paginazione
	startIdx := 0
	if cursor != "" {
		idx, err := strconv.Atoi(cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}
		startIdx = idx
	}

	endIdx := startIdx + limit
	if endIdx > len(posts) {
		endIdx = len(posts)
	}

	// Estrai slice di post per questa pagina
	postSlice := posts[startIdx:endIdx]

	// Converte in formato skeleton (solo URI)
	feed := make([]*bsky.FeedDefs_SkeletonFeedPost, len(postSlice))
	for i, sp := range postSlice {
		feed[i] = &bsky.FeedDefs_SkeletonFeedPost{
			Post: sp.AtURI,
		}
	}

	// Genera cursor per pagina successiva
	var nextCursor *string
	if endIdx < len(posts) {
		c := strconv.Itoa(endIdx)
		nextCursor = &c
	}

	// posts := make([]*bsky.FeedDefs_SkeletonFeedPost, len(feed))
	// t := "10"
	return &bsky.FeedGetFeedSkeleton_Output{
		Cursor: nextCursor,
		Feed:   feed,
	}, nil
}
