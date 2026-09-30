package userhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aakashloyar/elevate/problem/internal/application/ports/out"
	"github.com/aakashloyar/elevate/problem/internal/httpobserve"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) out.UserClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 5 * time.Second}}
}

func (c *Client) Exists(ctx context.Context, userID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/users/"+userID, nil)
	if err != nil {
		return err
	}
	resp, err := httpobserve.Do(c.http, req, "user.get")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("user %s was not found: %s", userID, resp.Status)
	}
	return nil
}

func (c *Client) GetUsernames(ctx context.Context, userIDs []string) (map[string]string, error) {
	result := make(map[string]string, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}
	body, err := json.Marshal(map[string][]string{"user_ids": userIDs})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/users/batch", strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpobserve.Do(c.http, req, "user.batch")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("user batch request failed: %s", resp.Status)
	}
	var payload struct {
		Users []struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"users"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	for _, user := range payload.Users {
		result[user.ID] = user.Username
	}
	return result, nil
}
