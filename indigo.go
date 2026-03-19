package main

import (
	"context"
	"fmt"

	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/xrpc"
)

func Indigo() {
	client := &xrpc.Client{Host: "https://api.bsky.app"}

	results, err := bsky.FeedSearchPosts(context.Background(), client, "", "", "", "", 25, "", "cat", "", "latest", nil, "", "")
	if err != nil {
		panic(err)
	}

	for _, post := range results.Posts {
		fmt.Printf("Author: %s\n", post.Author.Handle)
		fmt.Printf("Text: %s\n", post.Record.Val.(*bsky.FeedPost).Text)
		fmt.Println("---")
	}
}
