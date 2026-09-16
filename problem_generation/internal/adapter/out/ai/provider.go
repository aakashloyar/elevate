package ai

import (
	"fmt"
	"strings"

	"github.com/aakashloyar/elevate/problem_generation/internal/adapter/out/ai/gemini"
	"github.com/aakashloyar/elevate/problem_generation/internal/adapter/out/ai/groq"
	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
)

type Config struct {
	Provider string
	BaseURL  string
	APIKey   string
	Model    string
}

// NewClient selects a provider without coupling the generation workflow to a vendor SDK.
func NewClient(cfg Config) (out.AITextGenerator, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case "gemini", "google":
		baseURL := cfg.BaseURL
		if strings.Contains(baseURL, "api.groq.com") {
			baseURL = ""
		}
		return gemini.Config{BaseURL: baseURL, APIKey: cfg.APIKey, Model: cfg.Model}.NewClient(), nil
	case "groq", "grok":
		baseURL := cfg.BaseURL
		if strings.Contains(baseURL, "generativelanguage.googleapis.com") {
			baseURL = ""
		}
		return groq.Config{BaseURL: baseURL, APIKey: cfg.APIKey, Model: cfg.Model}.NewClient(), nil
	default:
		return nil, fmt.Errorf("unsupported AI provider %q (supported: gemini, groq)", cfg.Provider)
	}
}
