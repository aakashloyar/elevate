package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
)

type Client struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

type chatRequest struct {
	Model          string          `json:"model"`
	Messages       []message       `json:"messages"`
	Temperature    *float64        `json:"temperature,omitempty"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

func (c *Client) GenerateText(ctx context.Context, request out.GenerateTextRequest) (out.GenerateTextResponse, error) {
	if c.apiKey == "" {
		return out.GenerateTextResponse{}, errors.New("OpenAI-compatible API key is required")
	}
	if c.model == "" {
		return out.GenerateTextResponse{}, errors.New("OpenAI-compatible model is required")
	}

	payload := chatRequest{
		Model: c.model,
		Messages: []message{
			{Role: "system", Content: request.SystemPrompt},
			{Role: "user", Content: request.UserPrompt},
		},
		Temperature: request.Temperature,
	}
	if request.ResponseMimeType == "application/json" {
		payload.ResponseFormat = &responseFormat{Type: "json_object"}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return out.GenerateTextResponse{}, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return out.GenerateTextResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

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
		return out.GenerateTextResponse{}, fmt.Errorf("OpenAI-compatible request failed with status %d: %s", httpResp.StatusCode, string(responseBody))
	}

	var response chatResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return out.GenerateTextResponse{}, err
	}
	if len(response.Choices) == 0 || strings.TrimSpace(response.Choices[0].Message.Content) == "" {
		return out.GenerateTextResponse{}, errors.New("OpenAI-compatible response did not include text")
	}
	return out.GenerateTextResponse{Text: strings.TrimSpace(response.Choices[0].Message.Content)}, nil
}
