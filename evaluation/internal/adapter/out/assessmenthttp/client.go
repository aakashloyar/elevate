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
	"github.com/aakashloyar/elevate/evaluation/internal/httpobserve"
)

const defaultTimeout = 10 * time.Second

type Client struct {
	baseURL string
	http    *http.Client
}

type GetAssessmentProblemsResponse struct {
	ProblemIDs []string `json:"problem_ids"`
}

type GetAssessmentResponse struct {
	Title string `json:"title"`
}

type GetAssessmentsBatchResponse struct {
	Assessments []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"assessments"`
}

func NewClient(baseURL string) out.AssessmentClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: defaultTimeout}}
}

func (c *Client) GetAssessmentTitle(ctx context.Context, assessmentID string) (string, error) {
	var value GetAssessmentResponse
	if err := c.doGetRequest(ctx, "/assessments/"+assessmentID, &value); err != nil {
		return "", err
	}
	return value.Title, nil
}

func (c *Client) GetAssessmentTitles(ctx context.Context, assessmentIDs []string) (map[string]string, error) {
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
		return nil, fmt.Errorf("POST /assessments/batch returned %s", resp.Status)
	}
	var payload GetAssessmentsBatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	for _, assessment := range payload.Assessments {
		result[assessment.ID] = assessment.Title
	}
	return result, nil
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
	resp, err := httpobserve.Do(c.http, req, "assessment.get")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("GET %s returned %s", path, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(output)
}
