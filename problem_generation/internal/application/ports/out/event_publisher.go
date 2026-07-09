package out

import (
	"context"

	"github.com/aakashloyar/elevate/problem_generation/internal/domain"
)

type EventPublisher interface {
	PublishGenerationRequested(ctx context.Context, job domain.GenerationJob, prompt string) error
}
