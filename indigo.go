package main

import (
	"context"
	"fmt"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/xrpc"
)

func queryPost(handle, appPassword, query string, limit int64) {
	client := &xrpc.Client{Host: "https://bsky.social"}

	session, err := atproto.ServerCreateSession(context.Background(), client, &atproto.ServerCreateSession_Input{
		Identifier: handle,
		Password:   appPassword,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create session: %v", err))
	}

	authClient := &xrpc.Client{
		Host: "https://api.bsky.app",
		Auth: &xrpc.AuthInfo{
			AccessJwt:  session.AccessJwt,
			RefreshJwt: session.RefreshJwt,
		},
	}

	results, err := bsky.FeedSearchPosts(context.Background(), authClient, "", "", "", "it", limit, "", query, "", "top", nil, "", "")
	if err != nil {
		panic(err)
	}

	for _, post := range results.Posts {
		feedPost := post.Record.Val.(*bsky.FeedPost)

		fmt.Printf("\n=== POST ===\n")
		fmt.Printf("URI: %s\n", post.Uri)
		fmt.Printf("CID: %s\n", post.Cid)
		fmt.Printf("Author - Handle: %s\n", post.Author.Handle)
		fmt.Printf("Author - DID: %s\n", post.Author.Did)
		fmt.Printf("Author - DisplayName: %s\n", *post.Author.DisplayName)
		fmt.Printf("Text: %s\n", feedPost.Text)
		fmt.Printf("CreatedAt: %s\n", feedPost.CreatedAt)
		fmt.Printf("Languages: %v\n", feedPost.Langs)
		fmt.Printf("Tags: %v\n", feedPost.Tags)

		if feedPost.Embed != nil {
			if feedPost.Embed.EmbedImages != nil {
				fmt.Printf("Embed - Images: %d\n", len(feedPost.Embed.EmbedImages.Images))
				for i, img := range feedPost.Embed.EmbedImages.Images {
					fmt.Printf("  Image %d - Alt: %s\n", i, img.Alt)
				}
			}
			if feedPost.Embed.EmbedExternal != nil {
				fmt.Printf("Embed - External:\n")
				fmt.Printf("  URI: %s\n", feedPost.Embed.EmbedExternal.External.Uri)
				fmt.Printf("  Title: %s\n", feedPost.Embed.EmbedExternal.External.Title)
				fmt.Printf("  Description: %s\n", feedPost.Embed.EmbedExternal.External.Description)
			}
			if feedPost.Embed.EmbedRecord != nil {
				fmt.Printf("Embed - Record (quoted post):\n")
				fmt.Printf("  URI: %s\n", feedPost.Embed.EmbedRecord.Record.Uri)
			}
			if feedPost.Embed.EmbedVideo != nil {
				fmt.Printf("Embed - Video present\n")
			}
		}

		if feedPost.Reply != nil {
			fmt.Printf("Reply - Root URI: %s\n", feedPost.Reply.Root.Uri)
			fmt.Printf("Reply - Parent URI: %s\n", feedPost.Reply.Parent.Uri)
		}

		if feedPost.Facets != nil {
			fmt.Printf("Facets count: %d\n", len(feedPost.Facets))
		}

		fmt.Printf("IndexedAt: %s\n", post.IndexedAt)
		fmt.Printf("Labels: %v\n", post.Labels)
		fmt.Printf("LikeCount: %d\n", *post.LikeCount)
		fmt.Printf("RepostCount: %d\n", *post.RepostCount)
		fmt.Printf("ReplyCount: %d\n", *post.ReplyCount)
		fmt.Printf("QuoteCount: %d\n", *post.QuoteCount)

		if *post.ReplyCount > 0 {
			getReplies(authClient, post.Uri)
		}

		fmt.Printf("\n")
	}
}

func getReplies(authClient *xrpc.Client, postUri string) {
	thread, err := bsky.FeedGetPostThread(context.Background(), authClient, 3, 0, postUri)
	if err != nil {
		fmt.Printf("Error getting thread: %v\n", err)
		return
	}

	if thread.Thread == nil || thread.Thread.FeedDefs_ThreadViewPost == nil {
		return
	}

	viewPost := thread.Thread.FeedDefs_ThreadViewPost
	if viewPost.Replies == nil {
		return
	}

	fmt.Printf("=== REPLIES (%d) ===\n", len(viewPost.Replies))
	for i, reply := range viewPost.Replies {
		if reply.FeedDefs_ThreadViewPost != nil {
			replyPost := reply.FeedDefs_ThreadViewPost.Post
			replyText := ""
			if replyPost.Record != nil && replyPost.Record.Val != nil {
				if fp, ok := replyPost.Record.Val.(*bsky.FeedPost); ok {
					replyText = fp.Text
				}
			}
			fmt.Printf("  Reply %d - @%s: %s\n", i+1, replyPost.Author.Handle, replyText)
		}
	}
}
