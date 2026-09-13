package problemhttp

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) out.ProblemClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 5 * time.Second}}
}

func (c *Client) Exists(ctx context.Context, problemID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/problems/"+problemID, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("problem %s was not found: %s", problemID, resp.Status)
	}
	return nil
}

var _ out.ProblemClient = (*Client)(nil)
