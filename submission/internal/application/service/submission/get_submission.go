package submission

import (
	"context"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in/submission"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
)

type GetSubmissionService struct {
	submissionRepo out.SubmissionRepository
}

func NewGetSubmissionService(submissionRepo out.SubmissionRepository) in.GetSubmissionService {
	return &GetSubmissionService{submissionRepo: submissionRepo}
}

func (s *GetSubmissionService) Execute(ctx context.Context, input in.GetSubmissionInput) (in.GetSubmissionOutput, error) {
	submission, answers, err := s.submissionRepo.FindByID(input.SubmissionID)
	if err != nil {
		return in.GetSubmissionOutput{}, err
	}

	answerOutputs := make([]in.SubmissionAnswerOutput, 0, len(answers))
	for _, ans := range answers {
		answerOutputs = append(answerOutputs, in.SubmissionAnswerOutput{
			ProblemID:  ans.ProblemID,
			Answer:     ans.Answer,
			AnsweredAt: ans.AnsweredAt,
		})
	}

	return in.GetSubmissionOutput{
		ID:           submission.ID,
		AssessmentID: submission.AssessmentID,
		UserID:       submission.UserID,
		Status:       submission.Status,
		StartedAt:    submission.StartedAt,
		SubmittedAt:  submission.SubmittedAt,
		Answers:      answerOutputs,
	}, nil
}
