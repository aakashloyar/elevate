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

type SaveAnswerService struct {
	submissionRepo out.SubmissionRepository
	clock          out.Clock
}

func NewSaveAnswerService(submissionRepo out.SubmissionRepository, clock out.Clock) in.SaveAnswerService {
	return &SaveAnswerService{submissionRepo: submissionRepo, clock: clock}
}

func (s *SaveAnswerService) Execute(ctx context.Context, input in.SaveAnswerInput) error {
	if strings.TrimSpace(input.SubmissionID) == "" {
		return errors.New("submission id is required")
	}
	if strings.TrimSpace(input.ProblemID) == "" {
		return errors.New("problem id is required")
	}

	_, drafts, err := s.submissionRepo.FindAnswerSnapshots(input.SubmissionID)
	if err != nil {
		return err
	}
	snapshot, err := findProblemSnapshot(drafts, input.ProblemID)
	if err != nil {
		return err
	}
	answer := normalizeAnswer(input.Answer)
	if err := validateAnswer(snapshot, input.ProblemID, answer); err != nil {
		return err
	}

	return s.execute(input.SubmissionID, input.ProblemID, answer)
}

func (s *SaveAnswerService) execute(submissionID string, problemID string, answer []string) error {
	now := s.clock.Now()
	answerUpdatedAt := now
	saved, err := s.submissionRepo.SaveAnswer(domain.SubmissionAnswerDraft{
		SubmissionID:    submissionID,
		ProblemID:       problemID,
		Answer:          answer,
		AnswerUpdatedAt: &answerUpdatedAt,
	})
	if err != nil {
		return err
	}
	if !saved {
		return errors.New("submission must be in IN_PROGRESS state")
	}

	return nil
}

func normalizeAnswer(answer []string) []string {
	normalized := make([]string, 0, len(answer))
	for _, value := range answer {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}
	return normalized
}

func findProblemSnapshot(drafts []domain.SubmissionAnswerDraft, problemID string) (domain.SubmissionAnswerDraft, error) {
	for index := range drafts {
		if drafts[index].ProblemID == problemID {
			return drafts[index], nil
		}
	}
	return domain.SubmissionAnswerDraft{}, fmt.Errorf("problem %s is not part of submission", problemID)
}

func validateAnswer(snapshot domain.SubmissionAnswerDraft, problemID string, answer []string) error {
	if snapshot.ProblemType == domain.ProblemTypeNumerical {
		return nil
	}
	if len(answer) == 0 {
		return nil
	}
	validOptions := make(map[string]struct{}, len(snapshot.Options))
	for _, option := range snapshot.Options {
		validOptions[option.ID] = struct{}{}
	}
	for _, selectedOptionID := range answer {
		if _, ok := validOptions[selectedOptionID]; !ok {
			return fmt.Errorf("option %s is not valid for problem %s", selectedOptionID, problemID)
		}
	}
	return nil
}
