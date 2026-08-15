package problem

import (
	"context"
	"errors"
	"strings"

	in "github.com/aakashloyar/elevate/problem/internal/application/ports/in"
	"github.com/aakashloyar/elevate/problem/internal/application/ports/out"
	"github.com/aakashloyar/elevate/problem/internal/domain"
)

type UpdateProblemService struct {
	problemRepo out.ProblemRepository
	idGen       out.IDGenerator
	clock       out.Clock
}

func NewUpdateProblemService(problemRepo out.ProblemRepository, idGen out.IDGenerator, clock out.Clock) in.UpdateProblemService {
	return &UpdateProblemService{problemRepo: problemRepo, idGen: idGen, clock: clock}
}

func (s *UpdateProblemService) Execute(ctx context.Context, input in.UpdateProblemInput) error {
	statement := strings.TrimSpace(input.Statement)
	if statement == "" {
		return errors.New("statement is required")
	}

	createdBy := strings.TrimSpace(input.CreatedBy)
	if createdBy == "" {
		return errors.New("created by is required")
	}

	problem := domain.Problem{
		ID:         input.ProblemID,
		CreatedBy:  createdBy,
		Title:      strings.TrimSpace(input.Title),
		Statement:  statement,
		Type:       input.Type,
		Difficulty: input.Difficulty,
		SourceType: input.SourceType,
		UpdatedAt:  s.clock.Now(),
	}

	options := make([]domain.ProblemOption, 0, len(input.Options))
	for _, opt := range input.Options {
		options = append(options, domain.ProblemOption{
			ID:        s.idGen.NewID(),
			ProblemID: problem.ID,
			Text:      strings.TrimSpace(opt.Text),
			IsCorrect: opt.IsCorrect,
		})
	}

	tags := make([]domain.ProblemTag, 0, len(input.Tags))
	for _, tag := range input.Tags {
		tagName := strings.TrimSpace(tag)
		if tagName == "" {
			continue
		}
		tags = append(tags, domain.ProblemTag{ProblemID: problem.ID, Tag: tagName})
	}

	return s.problemRepo.Update(problem, options, tags)
}
