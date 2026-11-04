package documenthttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
	"github.com/aakashloyar/elevate/problem_generation/internal/httpobserve"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) out.DocumentClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 20 * time.Second}}
}

func (c *Client) GetDocument(ctx context.Context, documentID string) (out.DocumentSummary, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/documents/"+documentID, nil)
	if err != nil {
		return out.DocumentSummary{}, err
	}
	resp, err := httpobserve.Do(c.http, req, "document.get")
	if err != nil {
		return out.DocumentSummary{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return out.DocumentSummary{}, fmt.Errorf("document request failed: %s", resp.Status)
	}
	var document out.DocumentSummary
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return out.DocumentSummary{}, err
	}
	return document, nil
}

func (c *Client) GetContent(ctx context.Context, documentID string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/documents/"+documentID+"/content", nil)
	if err != nil {
		return "", err
	}
	resp, err := httpobserve.Do(c.http, req, "document.content")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("document content request failed: %s", resp.Status)
	}
	var payload struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	if strings.TrimSpace(payload.Content) == "" {
		return "", fmt.Errorf("document %s has no extracted content", documentID)
	}
	return payload.Content, nil
}

var _ out.DocumentClient = (*Client)(nil)
