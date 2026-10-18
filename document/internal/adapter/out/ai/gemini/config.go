package gemini

import "net/http"

type Config struct {
	BaseURL          string
	APIKey           string
	Model            string
	OutputDimensions int
	HTTPClient       *http.Client
}

func (c Config) withDefaults() Config {
	if c.BaseURL == "" {
		c.BaseURL = "https://generativelanguage.googleapis.com/v1beta"
	}
	if c.Model == "" {
		c.Model = "gemini-embedding-001"
	}
	if c.OutputDimensions == 0 {
		c.OutputDimensions = 768
	}
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}
	return c
}
