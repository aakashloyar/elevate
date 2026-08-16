package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *Client {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 5 * time.Second}}
}
func (c *Client) MarkingScheme(ctx context.Context, assessmentID string) (domain.MarkingScheme, error) {
	var value domain.MarkingScheme
	if err := c.get(ctx, "/assessments/"+assessmentID+"/marking-scheme", &value); err != nil {
		return domain.MarkingScheme{}, err
	}
	return value, nil
}
func (c *Client) ProblemIDs(ctx context.Context, assessmentID string) ([]string, error) {
	var value struct {
		ProblemIDs []string `json:"problem_ids"`
	}
	if err := c.get(ctx, "/assessments/"+assessmentID+"/problems", &value); err != nil {
		return nil, err
	}
	return value.ProblemIDs, nil
}
func (c *Client) Problem(ctx context.Context, problemID string) (domain.Problem, error) {
	var value domain.Problem
	if err := c.get(ctx, "/problems/"+problemID, &value); err != nil {
		return domain.Problem{}, err
	}
	return value, nil
}
func (c *Client) get(ctx context.Context, path string, output any) error {
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
