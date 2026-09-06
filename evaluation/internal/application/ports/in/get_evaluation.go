package evaluation

import (
	"context"

	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

type GetEvaluationService interface {
	Execute(ctx context.Context, submissionID string) (domain.Evaluation, error)
}
