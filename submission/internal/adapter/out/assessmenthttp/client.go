package assessmenthttp

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

func NewClient(baseURL string) out.AssessmentClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 5 * time.Second}}
}

func (c *Client) Exists(ctx context.Context, assessmentID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/assessments/"+assessmentID, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("assessment %s was not found: %s", assessmentID, resp.Status)
	}
	return nil
}

var _ out.AssessmentClient = (*Client)(nil)
