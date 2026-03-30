package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	_ "github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/xrpc"
)

func NewClient(handle, appPassword string) (*Client, error) {
	xrpcClient := &xrpc.Client{Host: "https://bsky.social"}

	session, err := atproto.ServerCreateSession(context.Background(), xrpcClient, &atproto.ServerCreateSession_Input{
		Identifier: handle,
		Password:   appPassword,
	})
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		client: &xrpc.Client{
			Host: "https://api.bsky.app",
			Auth: &xrpc.AuthInfo{
				AccessJwt:  session.AccessJwt,
				RefreshJwt: session.RefreshJwt,
			},
		},
	}, nil
}

func (c *Client) GetPost(ctx context.Context, handle string, key string) (Post, error) {
	atURI, err := c.GetAtUri(ctx, handle, key)
	if err != nil {
		return Post{}, err
	}

	result, err := bsky.FeedGetPosts(ctx, c.client, []string{atURI})
	if err != nil {
		return Post{}, err
	}

	if len(result.Posts) == 0 {
		return Post{}, fmt.Errorf("post not found: %s", atURI)
	}

	postView := result.Posts[0]
	record := postView.Record.Val.(*bsky.FeedPost)

	var labels []string
	for _, label := range postView.Labels {
		labels = append(labels, label.Val)
	}

	authorName := postView.Author.Handle
	if postView.Author.DisplayName != nil && *postView.Author.DisplayName != "" {
		authorName = *postView.Author.DisplayName
	}

	var replyRoot, replyParent string
	if record.Reply != nil {
		replyRoot = record.Reply.Root.Uri
		replyParent = record.Reply.Parent.Uri
	}

	url, err := buildURL(handle, key)
	if err != nil {
		fmt.Println("Error in buildURL")
		url = ""
	}

	return Post{
		URL:         url,
		AtURI:       atURI,
		Text:        record.Text,
		CreatedAt:   record.CreatedAt,
		Labels:      labels,
		Langs:       record.Langs,
		Tags:        record.Tags,
		Images:      extractImagesView(postView),
		Links:       extractLinksView(postView),
		Facets:      extractFacets(record),
		AuthorName:  authorName,
		ReplyRoot:   replyRoot,
		ReplyParent: replyParent,
	}, nil
}

// https://bsky.app/profile/{handle or DID}/post/{rkey}
func (c *Client) GetPostUrl(ctx context.Context, url string) (Post, error) {
	re := regexp.MustCompile(`^https://bsky\.app/profile/([^/]+)/post/([^/]+)$`)
	matches := re.FindStringSubmatch(url)

	handle := matches[1]
	key := matches[2]

	post, err := c.GetPost(ctx, handle, key)
	if err != nil {
		return Post{}, err
	}

	return post, nil
}

// at://did:plc:vwzwgnygau7ed7b7wt5ux7y2/app.bsky.feed.post/3k5nobkf2w72g
func (c *Client) GetAtUri(ctx context.Context, handle string, key string) (string, error) {
	result, err := atproto.IdentityResolveHandle(ctx, c.client, handle)
	if err != nil {
		return "", err
	}

	atURI := fmt.Sprintf("at://%s/app.bsky.feed.post/%s", result.Did, key)

	return atURI, nil
}

func extractImagesView(postView *bsky.FeedDefs_PostView) []PostImage {
	var images []PostImage
	if postView.Embed == nil || postView.Embed.EmbedImages_View == nil {
		return images
	}
	for _, img := range postView.Embed.EmbedImages_View.Images {
		images = append(images, PostImage{
			Alt:   img.Alt,
			Image: img.Fullsize,
		})
	}
	return images
}

func extractLinksView(postView *bsky.FeedDefs_PostView) []PostLink {
	var links []PostLink
	if postView.Embed == nil || postView.Embed.EmbedExternal_View == nil {
		return links
	}
	ext := postView.Embed.EmbedExternal_View.External
	thumb := ""
	if ext.Thumb != nil {
		thumb = *ext.Thumb
	}
	links = append(links, PostLink{
		Uri:         ext.Uri,
		Title:       ext.Title,
		Description: ext.Description,
		Thumb:       thumb,
	})
	return links
}

func extractImages(post *bsky.FeedPost) []PostImage {
	var images []PostImage
	if post.Embed == nil || post.Embed.EmbedImages == nil {
		return images
	}
	for _, img := range post.Embed.EmbedImages.Images {
		images = append(images, PostImage{
			Alt:   img.Alt,
			Image: img.Image.Ref.String(),
		})
	}
	return images
}

func extractLinks(post *bsky.FeedPost) []PostLink {
	var links []PostLink
	if post.Embed == nil || post.Embed.EmbedExternal == nil {
		return links
	}
	ext := post.Embed.EmbedExternal.External
	thumb := ""
	if ext.Thumb != nil {
		thumb = ext.Thumb.Ref.String()
	}
	links = append(links, PostLink{
		Uri:         ext.Uri,
		Title:       ext.Title,
		Description: ext.Description,
		Thumb:       thumb,
	})
	return links
}

func extractFacets(post *bsky.FeedPost) []PostFacet {
	var facets []PostFacet
	for _, facet := range post.Facets {
		for _, feature := range facet.Features {
			switch {
			case feature.RichtextFacet_Link != nil:
				facets = append(facets, PostFacet{
					Type:  "link",
					Value: feature.RichtextFacet_Link.Uri,
				})
			case feature.RichtextFacet_Mention != nil:
				facets = append(facets, PostFacet{
					Type:  "mention",
					Value: feature.RichtextFacet_Mention.Did,
				})
			case feature.RichtextFacet_Tag != nil:
				facets = append(facets, PostFacet{
					Type:  "tag",
					Value: feature.RichtextFacet_Tag.Tag,
				})
			}
		}
	}
	return facets
}

func SavePostsToJson(filename string, data interface{}) error {
	filename = fmt.Sprintf("%s_%s.json", filename, time.Now().Format("20060102150405"))
	bytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	if err := os.WriteFile(filename, bytes, 0o644); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}
	return nil
}

// "https://bsky.app/profile/pietrosalvatori.bsky.social/post/3mfqzcs2wck2y",
func buildURL(handle string, key string) (string, error) {
	if handle == "" || key == "" {
		return "", fmt.Errorf("Empty hanlde or key")
	}

	url := fmt.Sprintf("https://bsky.app/profile/%s/post/%s", handle, key)

	return url, nil
}
