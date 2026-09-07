package in

import (
	"context"

	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
)

type ProcessGenerationJobService interface {
	Execute(ctx context.Context, event out.GenerationRequestedEvent) error
}
