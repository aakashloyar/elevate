package submission

import (
	"context"
	"strings"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type ListSubmissionsService struct {
	repository out.SubmissionRepository
}

func NewListSubmissionsService(repository out.SubmissionRepository) in.ListSubmissionsService {
	return &ListSubmissionsService{repository: repository}
}

func (s *ListSubmissionsService) Execute(_ context.Context, input in.ListSubmissionsInput) (in.ListSubmissionsOutput, error) {
	var submissions []domain.Submission
	var err error
	if strings.TrimSpace(input.UserID) != "" {
		submissions, err = s.repository.ListByUserID(strings.TrimSpace(input.UserID))
	} else {
		submissions, err = s.repository.ListAll()
	}
	if err != nil {
		return in.ListSubmissionsOutput{}, err
	}
	return in.ListSubmissionsOutput{Submissions: submissions}, nil
}
