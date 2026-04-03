package bsky

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	cbg "github.com/whyrusleeping/cbor-gen"
)

var (
	_ cbg.CBORMarshaler   = (*SchwartzValuesRecord)(nil)
	_ cbg.CBORUnmarshaler = (*SchwartzValuesRecord)(nil)
)

func init() {
	lexutil.RegisterType("com.schwartz.values", &SchwartzValuesRecord{})
}

type SchwartzValuesRecord struct {
	LexiconTypeID string             `json:"$type" cborgen:"$type,const=com.schwartz.values"`
	Weights       map[string]float64 `json:"weights" cborgen:"weights"`
	UpdatedAt     string             `json:"updatedAt" cborgen:"updatedAt"`
}

func (r *SchwartzValuesRecord) MarshalCBOR(w io.Writer) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

func (r *SchwartzValuesRecord) UnmarshalCBOR(br io.Reader) error {
	var buf []byte
	if _, err := br.Read(buf); err != nil {
		return err
	}
	return json.Unmarshal(buf, r)
}

func NewClient(handle, appPassword string) (xrpc.Client, error) {
	xrpcClient := &xrpc.Client{Host: "https://bsky.social"}

	session, err := atproto.ServerCreateSession(context.Background(), xrpcClient,
		&atproto.ServerCreateSession_Input{
			Identifier: handle,
			Password:   appPassword,
		})
	if err != nil {
		return xrpc.Client{}, err
	}

	return xrpc.Client{
		Host: "https://bsky.social",
		Auth: &xrpc.AuthInfo{
			AccessJwt:  session.AccessJwt,
			RefreshJwt: session.RefreshJwt,
		},
	}, nil
}

func GetProfileInfo(ctx context.Context, client *xrpc.Client, handle string) (string, error) {
	profile, err := bsky.ActorGetProfile(ctx, client, handle)
	if err != nil {
		return "", fmt.Errorf("failed to get profile: %w", err)
	}
	if profile.DisplayName != nil {
		return *profile.DisplayName, nil
	}
	return profile.Handle, nil
}

func SaveWeights(ctx context.Context, client *xrpc.Client, handle string, weights map[string]float64) error {
	// Check if record already exists
	existingCID := ""
	resp, err := atproto.RepoGetRecord(ctx, client, "", "com.schwartz.values", handle, "main")
	if err == nil && resp.Cid != nil {
		// Record exists, get its CID for swap
		existingCID = *resp.Cid
		fmt.Println("Existing record found, CID:", existingCID)
	} else {
		// Record doesn't exist yet, create new
		fmt.Println("Creating new record (no existing CID)")
	}

	// Prepare new record
	record := &SchwartzValuesRecord{
		Weights:   weights,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	validate := false
	input := &atproto.RepoPutRecord_Input{
		Repo:       handle,
		Collection: "com.schwartz.values",
		Rkey:       "main",
		Record:     &lexutil.LexiconTypeDecoder{Val: record},
		Validate:   &validate,
	}

	// If record exists, specify CID for atomic swap
	if existingCID != "" {
		input.SwapRecord = &existingCID
	}

	// Save to PDS
	result, err := atproto.RepoPutRecord(ctx, client, input)
	if err != nil {
		return fmt.Errorf("failed to save weights to PDS: %w", err)
	}

	fmt.Println("Weights saved successfully!")
	fmt.Println("CID:", result.Cid)
	fmt.Printf("WEIGHTS: %v\n", weights)
	fmt.Println("UPDATEDAT:", record.UpdatedAt)

	return nil
}

func GetWeights(ctx context.Context, client *xrpc.Client, handle string) (map[string]float64, error) {
	resp, err := atproto.RepoGetRecord(ctx, client, "", "com.schwartz.values", handle, "main")
	if err != nil {
		return nil, fmt.Errorf("failed to get weights from PDS: %w", err)
	}

	if resp.Value == nil || resp.Value.Val == nil {
		return nil, fmt.Errorf("no weights found")
	}

	record, ok := resp.Value.Val.(*SchwartzValuesRecord)
	if !ok {
		return nil, fmt.Errorf("invalid record type")
	}

	return record.Weights, nil
}

func ResolveDID(ctx context.Context, client *xrpc.Client, handle string) (string, error) {
	identity, err := atproto.IdentityResolveHandle(ctx, client, handle)
	if err != nil {
		return "", err
	}
	return identity.Did, nil
}

func DeleteWeights(ctx context.Context, client *xrpc.Client, handle string) error {
	// Check if record already exists
	existingCID := ""
	resp, err := atproto.RepoGetRecord(ctx, client, "", "com.schwartz.values", handle, "main")
	if err == nil && resp.Cid != nil {
		// Record exists, get its CID for swap
		existingCID = *resp.Cid
		fmt.Println("Existing record found, CID:", existingCID)
	} else {
		// Record doesn't exist yet, create new
		fmt.Println("No weights where saved to PDS, returning...")
		return nil
	}

	input := &atproto.RepoDeleteRecord_Input{
		Repo:       handle,
		Collection: "com.schwartz.values",
		Rkey:       "main",
	}

	// If record exists, specify CID for atomic swap
	if existingCID != "" {
		input.SwapRecord = &existingCID
	}

	// Delete to PDS
	_, err = atproto.RepoDeleteRecord(ctx, client, input)
	if err != nil {
		return fmt.Errorf("failed to save weights to PDS: %w", err)
	}

	fmt.Println("Weights removed successfully!")

	return nil
}
