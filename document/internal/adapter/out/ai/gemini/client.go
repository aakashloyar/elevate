package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
	"github.com/aakashloyar/elevate/document/internal/observe"
)

type Client struct {
	baseURL          string
	apiKey           string
	model            string
	outputDimensions int
	httpClient       *http.Client
}

func NewClient(cfg Config) out.Embedder {
	cfg = cfg.withDefaults()
	return &Client{
		baseURL:          strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:           cfg.APIKey,
		model:            cfg.Model,
		outputDimensions: cfg.OutputDimensions,
		httpClient:       cfg.HTTPClient,
	}
}

type batchEmbedRequest struct {
	Requests []embedRequest `json:"requests"`
}

type embedRequest struct {
	Model              string             `json:"model"`
	Content            embedContent       `json:"content"`
	EmbedContentConfig embedContentConfig `json:"embedContentConfig,omitempty"`
}

type embedContent struct {
	Parts []embedPart `json:"parts"`
}

type embedPart struct {
	Text string `json:"text"`
}

type embedContentConfig struct {
	TaskType             string `json:"taskType,omitempty"`
	OutputDimensionality int    `json:"outputDimensionality,omitempty"`
}

type batchEmbedResponse struct {
	Embeddings []struct {
		Values []float32 `json:"values"`
	} `json:"embeddings"`
}

func (c *Client) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	defer observe.Dependency("gemini.batch_embed_contents")()
	if len(texts) == 0 {
		return [][]float32{}, nil
	}
	if strings.TrimSpace(c.apiKey) == "" {
		return nil, fmt.Errorf("Gemini API key is not configured")
	}

	requests := make([]embedRequest, len(texts))
	modelResource := "models/" + strings.TrimPrefix(c.model, "models/")
	for i, text := range texts {
		if strings.TrimSpace(text) == "" {
			return nil, fmt.Errorf("document chunk %d is empty", i)
		}
		requests[i] = embedRequest{
			Model:   modelResource,
			Content: embedContent{Parts: []embedPart{{Text: text}}},
			EmbedContentConfig: embedContentConfig{
				TaskType:             "RETRIEVAL_DOCUMENT",
				OutputDimensionality: c.outputDimensions,
			},
		}
	}

	payload, err := json.Marshal(batchEmbedRequest{Requests: requests})
	if err != nil {
		return nil, fmt.Errorf("marshal Gemini embedding request: %w", err)
	}

	url := fmt.Sprintf("%s/%s:batchEmbedContents", c.baseURL, modelResource)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create Gemini embedding request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call Gemini embeddings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("Gemini embeddings returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var result batchEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode Gemini embedding response: %w", err)
	}
	if len(result.Embeddings) != len(texts) {
		return nil, fmt.Errorf("Gemini returned %d embeddings for %d chunks", len(result.Embeddings), len(texts))
	}

	embeddings := make([][]float32, len(result.Embeddings))
	for i, embedding := range result.Embeddings {
		if len(embedding.Values) != c.outputDimensions {
			return nil, fmt.Errorf("Gemini embedding %d has %d dimensions, expected %d", i, len(embedding.Values), c.outputDimensions)
		}
		embeddings[i] = embedding.Values
	}
	return embeddings, nil
}
