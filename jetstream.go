package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	js "github.com/bluesky-social/jetstream/pkg/client"
	"github.com/bluesky-social/jetstream/pkg/models"
)

type Scheduler struct{}

type Record struct {
	Type  string   `json:"type"`
	Text  string   `json:"text"`
	Langs []string `json:"langs"`
}

func (s *Scheduler) printRecord(record *Record) {
	fmt.Printf("Text: %v\nLangs: %v\n", record.Text, record.Langs)
}

func (s *Scheduler) AddWork(ctx context.Context, repo string, evt *models.Event) error {
	time.Sleep(2000 * time.Millisecond)

	var record Record
	raw_json, err := evt.Commit.Record.MarshalJSON()
	if err != nil {
		log.Fatalf("Could not marshal record: %v", err)
	}
	err = json.Unmarshal(raw_json, &record)
	if err != nil {
		log.Fatalf("Could not unmarshal record: %v", err)
	}

	prettyJSON, _ := json.MarshalIndent(record, "", "  ")
	fmt.Printf("JSON: %s\n============\n", prettyJSON)

	s.printRecord(&record)

	return nil
}

func (s *Scheduler) Shutdown() {
	log.Println("Shutting down...")
	os.Exit(0)
}

func Jetstream() {
	handler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(handler)
	scheduler := &Scheduler{}
	config := js.ClientConfig{
		Compress:          true,
		WebsocketURL:      "wss://jetstream2.us-west.bsky.network/subscribe",
		WantedDids:        []string{},
		WantedCollections: []string{"app.bsky.feed.post"},
		MaxSize:           0,
		ExtraHeaders: map[string]string{
			"User-Agent": "jetstream-client/v0.0.1",
		},
	}
	ctx := context.Background()

	client, err := js.NewClient(&config, logger, scheduler)
	if err != nil {
		os.Exit(1)
	}

	err = client.ConnectAndRead(ctx, nil)
	if err != nil {
		log.Fatalf("Connection error: %v", err)
		os.Exit(2)
	}
}
