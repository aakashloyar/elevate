package problemhttp

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

func NewClient(baseURL string) out.ProblemClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: defaultTimeout}}
}

func (c *Client) GetProblemByID(ctx context.Context, problemID string) (domain.Problem, error) {
	var value domain.Problem
	if err := c.doGetRequest(ctx, "/problems/"+problemID, &value); err != nil {
		return domain.Problem{}, err
	}
	return value, nil
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
