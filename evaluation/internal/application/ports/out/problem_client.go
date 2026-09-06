package out

import (
	"context"

	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

type ProblemClient interface {
	GetProblemByID(ctx context.Context, problemID string) (domain.Problem, error)
}
