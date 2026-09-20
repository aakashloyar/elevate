package userhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aakashloyar/elevate/evaluation/internal/application/ports/out"
)

type Client struct {
	baseURL string
	http    *http.Client
}

type getUserResponse struct {
	Username string `json:"username"`
}

func NewClient(baseURL string) out.UserClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 5 * time.Second}}
}

func (c *Client) GetUserName(ctx context.Context, userID string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/users/"+userID, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("GET /users/%s returned %s", userID, resp.Status)
	}
	var value getUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&value); err != nil {
		return "", err
	}
	return value.Username, nil
}

var _ out.UserClient = (*Client)(nil)
