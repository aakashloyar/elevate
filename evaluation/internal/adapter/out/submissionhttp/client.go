package submissionhttp

import (
	"bytes"
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

type UpdateSubmissionStatusRequest struct {
	Status domain.SubmissionStatus `json:"status"`
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 5 * time.Second}}
}

func (c *Client) UpdateSubmissionStatus(ctx context.Context, submissionID string, status domain.SubmissionStatus) error {
	body, err := json.Marshal(UpdateSubmissionStatusRequest{Status: status})
	if err != nil {
		return err
	}

	path := "/submissions/" + submissionID + "/status"
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("PATCH %s returned %s", path, resp.Status)
	}
	return nil
}
