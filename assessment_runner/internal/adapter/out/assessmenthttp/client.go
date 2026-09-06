package assessmenthttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aakashloyar/elevate/assessment_runner/internal/application/ports/out"
)

type Client struct {
	baseURL string
	http    *http.Client
}

type GetAssessmentProblemsRequest struct {
	AssessmentID string
}

type GetAssessmentProblemsResponse struct {
	ProblemIDs []string `json:"problem_ids"`
}

func NewClient(baseURL string) out.AssessmentClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 5 * time.Second}}
}

func (c *Client) GetAssessmentProblemIDs(ctx context.Context, assessmentID string) ([]string, error) {
	request := GetAssessmentProblemsRequest{AssessmentID: assessmentID}
	var response GetAssessmentProblemsResponse
	if err := c.doGetRequest(ctx, "/assessments/"+url.PathEscape(request.AssessmentID)+"/problems", &response); err != nil {
		return nil, err
	}
	return response.ProblemIDs, nil
}

func (c *Client) doGetRequest(ctx context.Context, path string, output any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	response, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("GET %s returned %s", path, response.Status)
	}
	return json.NewDecoder(response.Body).Decode(output)
}
