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
	AssessmentID string                        `json:"assessment_id"`
	Answers      []GetSubmissionAnswerResponse `json:"answers"`
}

type GetSubmissionAnswerResponse struct {
	ProblemID string   `json:"problem_id"`
	Answer    []string `json:"answer"`
}

func NewClient(baseURL string) out.SubmissionClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 5 * time.Second}}
}

func (c *Client) GetAttemptDraft(ctx context.Context, attemptID string) (out.AttemptDraft, error) {
	request := GetSubmissionRequest{SubmissionID: attemptID}
	var response GetSubmissionResponse
	if err := c.doGetRequest(ctx, "/submissions/"+url.PathEscape(request.SubmissionID), &response); err != nil {
		return out.AttemptDraft{}, err
	}
	if response.AssessmentID == "" {
		return out.AttemptDraft{}, fmt.Errorf("submission response has no assessment_id")
	}
	answers := make([]out.DraftAnswer, 0, len(response.Answers))
	for _, answer := range response.Answers {
		answers = append(answers, out.DraftAnswer{
			ProblemID: answer.ProblemID,
			Answer:    answer.Answer,
		})
	}
	return out.AttemptDraft{AssessmentID: response.AssessmentID, Answers: answers}, nil
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
