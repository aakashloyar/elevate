package problemhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	in "github.com/aakashloyar/elevate/assessment_runner/internal/application/ports/in"
	"github.com/aakashloyar/elevate/assessment_runner/internal/application/ports/out"
)

type Client struct {
	baseURL string
	http    *http.Client
}

type GetProblemRequest struct {
	ProblemID string
}

type GetProblemResponse struct {
	ID         string                     `json:"id"`
	Title      string                     `json:"title"`
	Statement  string                     `json:"statement"`
	Type       string                     `json:"type"`
	Difficulty string                     `json:"difficulty"`
	Options    []GetProblemOptionResponse `json:"options"`
}

type GetProblemOptionResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func NewClient(baseURL string) out.ProblemClient {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), http: &http.Client{Timeout: 5 * time.Second}}
}

func (c *Client) GetProblemByID(ctx context.Context, problemID string) (in.ProblemView, error) {
	request := GetProblemRequest{ProblemID: problemID}
	var response GetProblemResponse
	if err := c.doGetRequest(ctx, "/problems/"+url.PathEscape(request.ProblemID), &response); err != nil {
		return in.ProblemView{}, err
	}
	options := make([]in.OptionView, 0, len(response.Options))
	for _, option := range response.Options {
		options = append(options, in.OptionView{ID: option.ID, Text: option.Text})
	}
	return in.ProblemView{ID: response.ID, Title: response.Title, Statement: response.Statement, Type: response.Type, Difficulty: response.Difficulty, Options: options}, nil
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
