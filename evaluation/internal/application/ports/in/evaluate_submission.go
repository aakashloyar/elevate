package evaluation

import (
	"context"

	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

type EvaluateSubmissionService interface {
	Execute(ctx context.Context, submission domain.SubmissionSubmitted) (domain.Evaluation, error)
}
