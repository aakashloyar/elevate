package problemhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type Client struct {
	baseURL string
	http    *http.Client
}

type problemSnapshotsRequest struct {
	ProblemIDs []string `json:"problem_ids"`
}

type problemSnapshotResponse struct {
	ProblemID   string                  `json:"problem_id"`
	ProblemType string                  `json:"problem_type"`
	Options     []problemOptionResponse `json:"options"`
}

type problemOptionResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func NewClient(baseURL string) out.ProblemClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 5 * time.Second}}
}

func (c *Client) GetProblemSnapshots(ctx context.Context, problemIDs []string) ([]out.ProblemSnapshot, error) {
	payload, err := json.Marshal(problemSnapshotsRequest{ProblemIDs: problemIDs})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/problems/batch", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("problem batch request failed: %s", resp.Status)
	}
	var response []problemSnapshotResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	snapshots := make([]out.ProblemSnapshot, 0, len(response))
	for _, item := range response {
		options := make([]out.Option, 0, len(item.Options))
		for _, option := range item.Options {
			options = append(options, out.Option{
				ID:   option.ID,
				Text: option.Text,
			})
		}
		snapshots = append(snapshots, out.ProblemSnapshot{
			ProblemID:   item.ProblemID,
			ProblemType: domain.ProblemType(item.ProblemType),
			Options:     options,
		})
	}
	return snapshots, nil
}

var _ out.ProblemClient = (*Client)(nil)
