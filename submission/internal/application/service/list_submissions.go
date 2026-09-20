package submission

import (
	"context"
	"errors"
	"strings"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
)

type ListSubmissionsService struct {
	repository out.SubmissionRepository
}

func NewListSubmissionsService(repository out.SubmissionRepository) in.ListSubmissionsService {
	return &ListSubmissionsService{repository: repository}
}

func (s *ListSubmissionsService) Execute(_ context.Context, input in.ListSubmissionsInput) (in.ListSubmissionsOutput, error) {
	if strings.TrimSpace(input.UserID) == "" {
		return in.ListSubmissionsOutput{}, errors.New("user id is required")
	}
	submissions, err := s.repository.ListByUserID(input.UserID)
	if err != nil {
		return in.ListSubmissionsOutput{}, err
	}
	return in.ListSubmissionsOutput{Submissions: submissions}, nil
}
