package assessmenthttp

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

func NewClient(baseURL string) out.AssessmentClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) Exists(ctx context.Context, assessmentID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/assessments/"+assessmentID, nil)
	if err != nil {
		return err
	}
	resp, err := httpobserve.Do(c.http, req, "assessment.get")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("assessment %s was not found: %s", assessmentID, resp.Status)
	}
	return nil
}

func (c *Client) GetTitles(ctx context.Context, assessmentIDs []string) (map[string]string, error) {
	result := make(map[string]string, len(assessmentIDs))
	if len(assessmentIDs) == 0 {
		return result, nil
	}
	body, err := json.Marshal(map[string][]string{"assessment_ids": assessmentIDs})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/assessments/batch", strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpobserve.Do(c.http, req, "assessment.batch")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("assessment batch request failed: %s", resp.Status)
	}
	var payload struct {
		Assessments []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"assessments"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	for _, assessment := range payload.Assessments {
		result[assessment.ID] = assessment.Title
	}
	return result, nil
}

var _ out.AssessmentClient = (*Client)(nil)
