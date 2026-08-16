package submissionhttp

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

type GetSubmissionRequest struct {
	SubmissionID string
}

type GetSubmissionResponse struct {
	AssessmentID string `json:"assessment_id"`
}

func NewClient(baseURL string) out.SubmissionGateway {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 5 * time.Second}}
}

func (c *Client) GetAttemptAssessmentID(ctx context.Context, attemptID string) (string, error) {
	request := GetSubmissionRequest{SubmissionID: attemptID}
	var response GetSubmissionResponse
	if err := c.doGetRequest(ctx, "/submissions/"+url.PathEscape(request.SubmissionID), &response); err != nil {
		return "", err
	}
	if response.AssessmentID == "" {
		return "", fmt.Errorf("submission response has no assessment_id")
	}
	return response.AssessmentID, nil
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
