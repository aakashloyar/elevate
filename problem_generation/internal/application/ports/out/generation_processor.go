package out

import (
	"context"

	"github.com/aakashloyar/elevate/problem_generation/internal/domain"
)

type GenerationProcessor interface {
	ProcessGeneration(ctx context.Context, job domain.GenerationJob) error
}
