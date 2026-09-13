package submission

import (
	"context"
	"errors"
	"strings"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type SaveAnswerBatchService struct {
	submissionRepo out.SubmissionRepository
	clock          out.Clock
}

func NewSaveAnswerBatchService(submissionRepo out.SubmissionRepository, clock out.Clock) in.SaveAnswerBatchService {
	return &SaveAnswerBatchService{submissionRepo: submissionRepo, clock: clock}
}

func (s *SaveAnswerBatchService) Execute(ctx context.Context, input in.SaveAnswerBatchInput) (in.SaveAnswerBatchOutput, error) {
	if strings.TrimSpace(input.SubmissionID) == "" {
		return in.SaveAnswerBatchOutput{}, errors.New("submission id is required")
	}
	if len(input.Answers) == 0 {
		return in.SaveAnswerBatchOutput{}, errors.New("answers are required")
	}

	_, snapshots, err := s.submissionRepo.FindAnswerSnapshots(input.SubmissionID)
	if err != nil {
		return in.SaveAnswerBatchOutput{}, err
	}
	snapshotsByProblem := make(map[string]domain.SubmissionAnswerDraft, len(snapshots))
	for _, snapshot := range snapshots {
		snapshotsByProblem[snapshot.ProblemID] = snapshot
	}

	now := s.clock.Now()
	answerUpdatedAt := now
	answers := make([]domain.SubmissionAnswerDraft, 0, len(input.Answers))
	validationErrors := make([]in.SaveAnswerBatchItemError, 0)
	for _, item := range input.Answers {
		problemID := strings.TrimSpace(item.ProblemID)
		if problemID == "" {
			validationErrors = append(validationErrors, in.SaveAnswerBatchItemError{
				ProblemID: item.ProblemID,
				Message:   "problem id is required",
			})
			continue
		}
		snapshot, ok := snapshotsByProblem[problemID]
		if !ok {
			validationErrors = append(validationErrors, in.SaveAnswerBatchItemError{
				ProblemID: problemID,
				Message:   "problem " + problemID + " is not part of submission",
			})
			continue
		}
		answer := normalizeAnswer(item.Answer)
		if err := validateAnswer(snapshot, problemID, answer); err != nil {
			validationErrors = append(validationErrors, in.SaveAnswerBatchItemError{
				ProblemID: problemID,
				Message:   err.Error(),
			})
			continue
		}
		answers = append(answers, domain.SubmissionAnswerDraft{
			SubmissionID:    input.SubmissionID,
			ProblemID:       problemID,
			Answer:          answer,
			AnswerUpdatedAt: &answerUpdatedAt,
		})
	}

	if len(answers) == 0 {
		return in.SaveAnswerBatchOutput{SavedCount: 0, Errors: validationErrors}, nil
	}

	saved, err := s.submissionRepo.SaveAnswers(answers)
	if err != nil {
		return in.SaveAnswerBatchOutput{}, err
	}
	if !saved {
		return in.SaveAnswerBatchOutput{}, errors.New("submission must be in IN_PROGRESS state")
	}

	return in.SaveAnswerBatchOutput{SavedCount: len(answers), Errors: validationErrors}, nil
}
