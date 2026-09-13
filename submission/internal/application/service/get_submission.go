package submission

import (
	"context"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
)

type GetSubmissionService struct {
	submissionRepo out.SubmissionRepository
}

func NewGetSubmissionService(submissionRepo out.SubmissionRepository) in.GetSubmissionService {
	return &GetSubmissionService{submissionRepo: submissionRepo}
}

func (s *GetSubmissionService) Execute(ctx context.Context, input in.GetSubmissionInput) (in.GetSubmissionOutput, error) {
	submission, drafts, err := s.submissionRepo.FindByID(input.SubmissionID)
	if err != nil {
		return in.GetSubmissionOutput{}, err
	}

	problemOutputs := make([]in.SubmissionProblemOutput, 0, len(drafts))
	for _, draft := range drafts {
		optionIDs := make([]string, 0, len(draft.Options))
		optionTexts := make([]string, 0, len(draft.Options))
		for _, option := range draft.Options {
			optionIDs = append(optionIDs, option.ID)
			optionTexts = append(optionTexts, option.Text)
		}
		problemOutputs = append(problemOutputs, in.SubmissionProblemOutput{
			ProblemID:       draft.ProblemID,
			ProblemType:     draft.ProblemType,
			OptionIDs:       optionIDs,
			OptionTexts:     optionTexts,
			AnswerUpdatedAt: draft.AnswerUpdatedAt,
		})
	}
	answerOutputs := make([]in.SubmissionAnswerOutput, 0, len(drafts))
	for _, draft := range drafts {
		if draft.AnswerUpdatedAt != nil {
			answerOutputs = append(answerOutputs, in.SubmissionAnswerOutput{ProblemID: draft.ProblemID, Answer: draft.Answer, CreatedAt: *draft.AnswerUpdatedAt, UpdatedAt: *draft.AnswerUpdatedAt})
		}
	}

	return in.GetSubmissionOutput{
		ID:           submission.ID,
		AssessmentID: submission.AssessmentID,
		UserID:       submission.UserID,
		Status:       submission.Status,
		StartedAt:    submission.StartedAt,
		ExpiresAt:    submission.ExpiresAt,
		SubmittedAt:  submission.SubmittedAt,
		CreatedAt:    submission.CreatedAt,
		UpdatedAt:    submission.UpdatedAt,
		Answers:      answerOutputs,
		Problems:     problemOutputs,
	}, nil
}
