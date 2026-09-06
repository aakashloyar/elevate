package evaluation

import (
	"context"
	"errors"

	"github.com/aakashloyar/elevate/evaluation/internal/application/ports/out"
	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

var ErrEvaluationNotFound = errors.New("evaluation not found")

type GetEvaluationService struct {
	repository out.Repository
}

func NewGetEvaluationService(repository out.Repository) *GetEvaluationService {
	return &GetEvaluationService{repository: repository}
}

func (s *GetEvaluationService) Execute(ctx context.Context, submissionID string) (domain.Evaluation, error) {
	return s.repository.FindBySubmissionID(ctx, submissionID)
}
