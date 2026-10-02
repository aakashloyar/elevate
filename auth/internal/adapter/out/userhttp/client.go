package userhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aakashloyar/elevate/auth/internal/httpobserve"
)

type Client struct {
	base string
	http *http.Client
}

func New(base string) *Client {
	return &Client{base: strings.TrimRight(base, "/"), http: http.DefaultClient}
}

func (c *Client) Create(ctx context.Context, username, email string) (string, error) {
	body, _ := json.Marshal(map[string]string{"username": username, "email": email})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/users", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpobserve.Do(c.http, req, "user")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("user service returned %s", resp.Status)
	}
	var out struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.UserID, nil
}
