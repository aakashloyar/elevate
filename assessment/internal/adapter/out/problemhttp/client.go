package problemhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
)

const defaultTimeout = 15 * time.Second

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string) out.ProblemClient {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: defaultTimeout},
	}
}

func (c *Client) CreateProblem(ctx context.Context, input out.CreateProblemInput) (out.CreateProblemOutput, error) {
	payload, err := json.Marshal(input)
	if err != nil {
		return out.CreateProblemOutput{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/problems", bytes.NewReader(payload))
	if err != nil {
		return out.CreateProblemOutput{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return out.CreateProblemOutput{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		message := strings.TrimSpace(string(body))
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
