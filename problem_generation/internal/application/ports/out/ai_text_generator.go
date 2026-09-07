package out

import "context"

type GenerateTextRequest struct {
	SystemPrompt     string
	UserPrompt       string
	ResponseMimeType string
	Temperature      *float64
}

type GenerateTextResponse struct {
	Text string
}

type AITextGenerator interface {
	GenerateText(ctx context.Context, request GenerateTextRequest) (GenerateTextResponse, error)
}
