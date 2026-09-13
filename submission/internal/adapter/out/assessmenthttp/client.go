package assessmenthttp

import (
	"context"
	"encoding/json"
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

type assessmentProblemsResponse struct {
	ProblemIDs []string `json:"problem_ids"`
}

func NewClient(baseURL string) out.AssessmentClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 5 * time.Second}}
}

func (c *Client) GetProblemIDs(ctx context.Context, assessmentID string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/assessments/"+assessmentID+"/problems", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("assessment %s was not found: %s", assessmentID, resp.Status)
	}
	var value assessmentProblemsResponse
	if err := json.NewDecoder(resp.Body).Decode(&value); err != nil {
		return nil, err
	}
	return value.ProblemIDs, nil
}

var _ out.AssessmentClient = (*Client)(nil)
