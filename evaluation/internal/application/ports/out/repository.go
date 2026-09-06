package out

import (
	"context"

	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

type Repository interface {
	Save(ctx context.Context, evaluation domain.Evaluation) error
	FindBySubmissionID(ctx context.Context, submissionID string) (domain.Evaluation, error)
}
