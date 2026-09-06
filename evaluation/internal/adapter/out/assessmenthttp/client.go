package assessmenthttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aakashloyar/elevate/evaluation/internal/application/ports/out"
	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

const defaultTimeout = 5 * time.Second

type Client struct {
	baseURL string
	http    *http.Client
}

type GetAssessmentProblemsResponse struct {
	ProblemIDs []string `json:"problem_ids"`
}

func NewClient(baseURL string) out.AssessmentClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: defaultTimeout}}
}

func (c *Client) GetAssessmentMarkingScheme(ctx context.Context, assessmentID string) (domain.MarkingScheme, error) {
	var value domain.MarkingScheme
	if err := c.doGetRequest(ctx, "/assessments/"+assessmentID+"/marking-scheme", &value); err != nil {
		return domain.MarkingScheme{}, err
	}
	return value, nil
}

func (c *Client) GetAssessmentProblemIDs(ctx context.Context, assessmentID string) ([]string, error) {
	var value GetAssessmentProblemsResponse
	if err := c.doGetRequest(ctx, "/assessments/"+assessmentID+"/problems", &value); err != nil {
		return nil, err
	}
	return value.ProblemIDs, nil
}

func (c *Client) doGetRequest(ctx context.Context, path string, output any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("GET %s returned %s", path, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(output)
}
