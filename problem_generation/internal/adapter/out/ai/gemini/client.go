package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
)

type Client struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

type generateContentRequest struct {
	SystemInstruction content          `json:"system_instruction,omitempty"`
	Contents          []content        `json:"contents"`
	GenerationConfig  generationConfig `json:"generationConfig,omitempty"`
}

type generationConfig struct {
	ResponseMimeType string   `json:"responseMimeType,omitempty"`
	Temperature      *float64 `json:"temperature,omitempty"`
}

type content struct {
	Role  string `json:"role,omitempty"`
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

type generateContentResponse struct {
	Candidates []candidate `json:"candidates"`
}

type candidate struct {
	Content content `json:"content"`
}

func (c *Client) GenerateText(ctx context.Context, request out.GenerateTextRequest) (out.GenerateTextResponse, error) {
	if c.apiKey == "" {
		return out.GenerateTextResponse{}, errors.New("gemini api key is required")
	}
	if c.model == "" {
		return out.GenerateTextResponse{}, errors.New("gemini model is required")
	}

	payload := generateContentRequest{
		SystemInstruction: content{
			Parts: []part{{Text: request.SystemPrompt}},
		},
		Contents: []content{
			{
				Role:  "user",
				Parts: []part{{Text: request.UserPrompt}},
			},
		},
		GenerationConfig: generationConfig{
			ResponseMimeType: request.ResponseMimeType,
			Temperature:      request.Temperature,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return out.GenerateTextResponse{}, err
	}

	endpoint := fmt.Sprintf("%s/v1beta/models/%s:generateContent", c.baseURL, url.PathEscape(strings.TrimPrefix(c.model, "models/")))
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return out.GenerateTextResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", c.apiKey)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return out.GenerateTextResponse{}, err
	}
	defer httpResp.Body.Close()

	responseBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return out.GenerateTextResponse{}, err
	}

	if httpResp.StatusCode < http.StatusOK || httpResp.StatusCode >= http.StatusMultipleChoices {
		return out.GenerateTextResponse{}, fmt.Errorf("gemini request failed with status %d: %s", httpResp.StatusCode, string(responseBody))
	}

	var response generateContentResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return out.GenerateTextResponse{}, err
	}

	text := strings.TrimSpace(response.Text())
	if text == "" {
		return out.GenerateTextResponse{}, errors.New("gemini response did not include text")
	}

	return out.GenerateTextResponse{Text: text}, nil
}

func (r generateContentResponse) Text() string {
	if len(r.Candidates) == 0 {
		return ""
	}

	parts := r.Candidates[0].Content.Parts
	text := make([]string, 0, len(parts))
	for _, part := range parts {
		if strings.TrimSpace(part.Text) != "" {
			text = append(text, part.Text)
		}
	}
	return strings.Join(text, "")
}
