package submission

import (
	"context"
	"errors"
	"fmt"
	"strings"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type UpdateSubmissionStatusService struct {
	submissionRepo out.SubmissionRepository
}

func NewUpdateSubmissionStatusService(submissionRepo out.SubmissionRepository) in.UpdateSubmissionStatusService {
	return &UpdateSubmissionStatusService{submissionRepo: submissionRepo}
}

func (s *UpdateSubmissionStatusService) Execute(_ context.Context, input in.UpdateSubmissionStatusInput) (in.UpdateSubmissionStatusOutput, error) {
	if strings.TrimSpace(input.SubmissionID) == "" {
		return in.UpdateSubmissionStatusOutput{}, errors.New("submission id is required")
	}
	if input.Status == "" {
		return in.UpdateSubmissionStatusOutput{}, errors.New("submission status is required")
	}

	currentStatus, _, err := s.submissionRepo.FindStatus(input.SubmissionID)
	if err != nil {
		return in.UpdateSubmissionStatusOutput{}, err
	}
	if !canUpdateSubmissionStatus(currentStatus, input.Status) {
		return in.UpdateSubmissionStatusOutput{}, fmt.Errorf("cannot update submission status from %q to %q", currentStatus, input.Status)
	}

	if err := s.submissionRepo.UpdateStatus(input.SubmissionID, input.Status); err != nil {
		return in.UpdateSubmissionStatusOutput{}, err
	}
	return in.UpdateSubmissionStatusOutput{SubmissionID: input.SubmissionID, Status: input.Status}, nil
}

func canUpdateSubmissionStatus(currentStatus, nextStatus domain.SubmissionStatus) bool {
	if currentStatus == nextStatus {
		return true
	}

	allowedTransitions := map[domain.SubmissionStatus][]domain.SubmissionStatus{
		domain.SubmissionStatusSubmitted: {
			domain.SubmissionStatusUnderEvaluation,
		},
		domain.SubmissionStatusUnderEvaluation: {
			domain.SubmissionStatusEvaluated,
			domain.SubmissionStatusEvaluationFailed,
		},
		domain.SubmissionStatusEvaluationFailed: {
			domain.SubmissionStatusUnderEvaluation,
		},
	}

	for _, allowedStatus := range allowedTransitions[currentStatus] {
		if nextStatus == allowedStatus {
			return true
		}
	}
	return false
}
