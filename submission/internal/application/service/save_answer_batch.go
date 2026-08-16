package submission

import (
	"context"
	"errors"
	"strings"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
)

type SaveAnswerBatchService struct {
	saveAnswerService in.SaveAnswerService
}

func NewSaveAnswerBatchService(saveAnswerService in.SaveAnswerService) in.SaveAnswerBatchService {
	return &SaveAnswerBatchService{saveAnswerService: saveAnswerService}
}

func (s *SaveAnswerBatchService) Execute(ctx context.Context, input in.SaveAnswerBatchInput) error {
	if strings.TrimSpace(input.SubmissionID) == "" {
		return errors.New("submission id is required")
	}
	if len(input.Answers) == 0 {
		return errors.New("answers are required")
	}

	for _, item := range input.Answers {
		if err := s.saveAnswerService.Execute(ctx, in.SaveAnswerInput{
			SubmissionID: input.SubmissionID,
			ProblemID:    item.ProblemID,
			Answer:       item.Answer,
		}); err != nil {
			return err
		}
	}

	return nil
}
