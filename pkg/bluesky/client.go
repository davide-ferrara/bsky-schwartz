package bluesky

import (
	"context"
	"fmt"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/xrpc"

	"bsky-schwarz/pkg/scorer"
)

type Client struct {
	client *xrpc.Client
}

func NewClient(handle, appPassword string) Client {
	client := &xrpc.Client{Host: "https://bsky.social"}

	session, err := atproto.ServerCreateSession(context.Background(), client, &atproto.ServerCreateSession_Input{
		Identifier: handle,
		Password:   appPassword,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create session: %v", err))
	}

	return Client{
		client: &xrpc.Client{
			Host: "https://api.bsky.app",
			Auth: &xrpc.AuthInfo{
				AccessJwt:  session.AccessJwt,
				RefreshJwt: session.RefreshJwt,
			},
		},
	}
}

func (c Client) QueryPosts(query string, limit int64) []scorer.FeedItem {
	results, err := bsky.FeedSearchPosts(context.Background(), c.client, "", "", "", "it", limit, "", query, "", "top", nil, "", "")
	if err != nil {
		panic(err)
	}

	items := make([]scorer.FeedItem, 0, len(results.Posts))
	for _, post := range results.Posts {
		feedPost := post.Record.Val.(*bsky.FeedPost)

		item := scorer.FeedItem{
			URI:  post.Uri,
			Text: feedPost.Text,
		}

		if post.ReplyCount != nil && *post.ReplyCount > 0 {
			item.Replies = c.GetReplies(post.Uri)
		}

		if feedPost.Embed != nil {
			if feedPost.Embed.EmbedImages != nil {
				for _, img := range feedPost.Embed.EmbedImages.Images {
					item.Images = append(item.Images, scorer.EmbedImage{
						URI: fmt.Sprintf("https://cdn.bsky.app/img/feed_fullsize/plain/%s/%s@jpeg", post.Author.Did, img.Image.Ref.String()),
						Alt: img.Alt,
					})
				}
			}
			if feedPost.Embed.EmbedExternal != nil {
				item.External = &scorer.EmbedExternal{
					URI:         feedPost.Embed.EmbedExternal.External.Uri,
					Title:       feedPost.Embed.EmbedExternal.External.Title,
					Description: feedPost.Embed.EmbedExternal.External.Description,
				}
			}
			if feedPost.Embed.EmbedRecord != nil {
				item.QuotedPostURI = feedPost.Embed.EmbedRecord.Record.Uri
			}
		}

		items = append(items, item)
	}

	return items
}

func (c Client) GetPostsUri(query string, limit int64) []string {
	res, err := bsky.FeedSearchPosts(context.Background(), c.client, "", "", "", "it", limit, "", query, "", "top", nil, "", "")
	if err != nil {
		panic(err)
	}
	var uris []string
	for _, post := range res.Posts {
		uris = append(uris, post.Uri)
	}

	return uris
}

func (c Client) LogPost(post *bsky.FeedDefs_PostView, feedPost *bsky.FeedPost) {
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
		c.GetReplies(post.Uri)
	}

	fmt.Printf("\n")
}

func (c Client) GetReplies(postUri string) []string {
	thread, err := bsky.FeedGetPostThread(context.Background(), c.client, 3, 0, postUri)
	if err != nil {
		fmt.Printf("Error getting thread: %v\n", err)
		return nil
	}

	if thread.Thread == nil || thread.Thread.FeedDefs_ThreadViewPost == nil {
		return nil
	}

	viewPost := thread.Thread.FeedDefs_ThreadViewPost
	if viewPost.Replies == nil {
		return nil
	}

	var replies []string
	for _, reply := range viewPost.Replies {
		if reply.FeedDefs_ThreadViewPost != nil {
			replyPost := reply.FeedDefs_ThreadViewPost.Post
			if replyPost.Record != nil && replyPost.Record.Val != nil {
				if fp, ok := replyPost.Record.Val.(*bsky.FeedPost); ok {
					replies = append(replies, fp.Text)
				}
			}
		}
	}
	return replies
}

func (c Client) GetPost(postUri string) *scorer.FeedItem {
	thread, err := bsky.FeedGetPostThread(context.Background(), c.client, 1, 0, postUri)
	if err != nil {
		fmt.Printf("Error getting post: %v\n", err)
		return nil
	}

	if thread.Thread == nil || thread.Thread.FeedDefs_ThreadViewPost == nil {
		return nil
	}

	viewPost := thread.Thread.FeedDefs_ThreadViewPost.Post
	feedPost, ok := viewPost.Record.Val.(*bsky.FeedPost)
	if !ok {
		return nil
	}

	item := &scorer.FeedItem{
		URI:  viewPost.Uri,
		Text: feedPost.Text,
	}

	return item
}
