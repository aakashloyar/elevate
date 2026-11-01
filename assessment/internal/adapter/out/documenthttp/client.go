package documenthttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
	"github.com/aakashloyar/elevate/assessment/internal/httpobserve"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) out.DocumentClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: http.DefaultClient}
}

func (c *Client) GetDocuments(ctx context.Context, documentIDs []string) ([]out.DocumentSummary, error) {
	body, err := json.Marshal(map[string][]string{"document_ids": documentIDs})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/documents/batch", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpobserve.Do(c.http, req, "document.batch")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("document batch request failed: %s", resp.Status)
	}
	var payload struct {
		Documents []out.DocumentSummary `json:"documents"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return payload.Documents, nil
}
