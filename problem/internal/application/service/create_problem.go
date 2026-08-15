package problem

import (
	"context"
	"errors"
	"strings"

	in "github.com/aakashloyar/elevate/problem/internal/application/ports/in"
	"github.com/aakashloyar/elevate/problem/internal/application/ports/out"
	"github.com/aakashloyar/elevate/problem/internal/domain"
)

type CreateProblemService struct {
	problemRepo out.ProblemRepository
	idGen       out.IDGenerator
	clock       out.Clock
}

func NewCreateProblemService(problemRepo out.ProblemRepository, idGen out.IDGenerator, clock out.Clock) in.CreateProblemService {
	return &CreateProblemService{problemRepo: problemRepo, idGen: idGen, clock: clock}
}

func (s *CreateProblemService) Execute(ctx context.Context, input in.CreateProblemInput) (in.CreateProblemOutput, error) {
	statement := strings.TrimSpace(input.Statement)
	if statement == "" {
		return in.CreateProblemOutput{}, errors.New("statement is required")
	}

	createdBy := strings.TrimSpace(input.CreatedBy)
	if createdBy == "" {
		return in.CreateProblemOutput{}, errors.New("created by is required")
	}

	problem := domain.Problem{
		ID:         s.idGen.NewID(),
		CreatedBy:  createdBy,
		Title:      strings.TrimSpace(input.Title),
		Statement:  statement,
		Type:       input.Type,
		Difficulty: input.Difficulty,
		SourceType: input.SourceType,
		CreatedAt:  s.clock.Now(),
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

	if err := s.problemRepo.Save(problem, options, tags); err != nil {
		return in.CreateProblemOutput{}, err
	}

	return in.CreateProblemOutput{ProblemID: problem.ID}, nil
}
