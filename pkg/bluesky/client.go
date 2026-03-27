package bluesky

import (
	"context"
	"fmt"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/xrpc"

	"bsky-schwarz/pkg/logger"
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
		logger.Error("failed to create bluesky session", "error", err)
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
	logger.Info("bluesky search started", "query", query, "limit", limit)
	start := time.Now()

	results, err := bsky.FeedSearchPosts(context.Background(), c.client, "", "", "", "it", limit, "", query, "", "top", nil, "", "")
	if err != nil {
		logger.Error("bluesky search failed", "error", err, "query", query)
		panic(err)
	}

	logger.Info("bluesky search completed",
		"query", query,
		"found_posts", len(results.Posts),
		"duration_ms", time.Since(start).Milliseconds(),
	)

	items := make([]scorer.FeedItem, 0, len(results.Posts))
	for i, post := range results.Posts {
		postStart := time.Now()

		feedPost := post.Record.Val.(*bsky.FeedPost)

		item := scorer.FeedItem{
			URI:  post.Uri,
			Text: feedPost.Text,
		}

		// Replies disabled: comments can distort post analysis
		// Negative comments on positive posts shouldn't affect post ranking
		// if post.ReplyCount != nil && *post.ReplyCount > 0 {
		// 	item.Replies = c.GetReplies(post.Uri)
		// }

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

		logger.Debug("post processed",
			"index", i,
			"post_uri", post.Uri,
			"has_replies", item.Replies != nil,
			"duration_ms", time.Since(postStart).Milliseconds(),
		)
	}

	return items
}

func (c Client) GetPostsUri(query string, limit int64) []string {
	logger.Info("bluesky uri search started", "query", query, "limit", limit)
	start := time.Now()

	res, err := bsky.FeedSearchPosts(context.Background(), c.client, "", "", "", "it", limit, "", query, "", "top", nil, "", "")
	if err != nil {
		logger.Error("bluesky uri search failed", "error", err, "query", query)
		panic(err)
	}

	var uris []string
	for _, post := range res.Posts {
		uris = append(uris, post.Uri)
	}

	logger.Info("bluesky uri search completed",
		"query", query,
		"found_uris", len(uris),
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return uris
}

func (c Client) GetReplies(postUri string) []string {
	logger.Debug("fetching replies", "post_uri", postUri)
	start := time.Now()

	thread, err := bsky.FeedGetPostThread(context.Background(), c.client, 3, 0, postUri)
	if err != nil {
		logger.Warn("failed to get thread", "error", err, "post_uri", postUri)
		return nil
	}

	if thread.Thread == nil || thread.Thread.FeedDefs_ThreadViewPost == nil {
		logger.Debug("no thread found", "post_uri", postUri)
		return nil
	}

	viewPost := thread.Thread.FeedDefs_ThreadViewPost
	if viewPost.Replies == nil {
		logger.Debug("no replies found", "post_uri", postUri)
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

	logger.Debug("replies fetched",
		"post_uri", postUri,
		"count", len(replies),
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return replies
}

func (c Client) GetPost(postUri string) *scorer.FeedItem {
	logger.Info("fetching single post", "post_uri", postUri)
	start := time.Now()

	thread, err := bsky.FeedGetPostThread(context.Background(), c.client, 1, 0, postUri)
	if err != nil {
		logger.Error("failed to get post", "error", err, "post_uri", postUri)
		return nil
	}

	if thread.Thread == nil || thread.Thread.FeedDefs_ThreadViewPost == nil {
		logger.Warn("post not found", "post_uri", postUri)
		return nil
	}

	viewPost := thread.Thread.FeedDefs_ThreadViewPost.Post
	feedPost, ok := viewPost.Record.Val.(*bsky.FeedPost)
	if !ok {
		logger.Warn("failed to parse post", "post_uri", postUri)
		return nil
	}

	item := &scorer.FeedItem{
		URI:  viewPost.Uri,
		Text: feedPost.Text,
	}

	logger.Info("post fetched successfully",
		"post_uri", postUri,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return item
}
