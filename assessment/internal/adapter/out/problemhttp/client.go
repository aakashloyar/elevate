package problemhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) out.ProblemClient {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: http.DefaultClient,
	}
}

func (c *Client) CreateProblem(ctx context.Context, input out.CreateProblemInput) (out.CreateProblemOutput, error) {
	if c.baseURL == "" {
		return out.CreateProblemOutput{}, errors.New("problem service url is required")
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return out.CreateProblemOutput{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/problems", bytes.NewReader(payload))
	if err != nil {
		return out.CreateProblemOutput{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return out.CreateProblemOutput{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(resp.Body)
		message := strings.TrimSpace(string(body))
		if message == "" {
			message = http.StatusText(resp.StatusCode)
		}
		return out.CreateProblemOutput{}, &out.ProblemClientError{StatusCode: resp.StatusCode, Message: message}
	}

	var output out.CreateProblemOutput
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return out.CreateProblemOutput{}, err
	}
	if strings.TrimSpace(output.ProblemID) == "" {
		return out.CreateProblemOutput{}, &out.ProblemClientError{StatusCode: resp.StatusCode, Message: "missing problem_id in response"}
	}
	return output, nil
}

func (c *Client) Exists(ctx context.Context, problemID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/problems/"+problemID, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("problem %s was not found: %s", problemID, resp.Status)
	}
	return nil
}
