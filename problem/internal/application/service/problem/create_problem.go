package problem

import (
	"context"
	"errors"
	"strings"

	in "github.com/aakashloyar/elevate/problem/internal/application/ports/in/problem"
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
		Type:       strings.ToUpper(strings.TrimSpace(input.Type)),
		Difficulty: strings.ToUpper(strings.TrimSpace(input.Difficulty)),
		SourceType: strings.ToUpper(strings.TrimSpace(input.SourceType)),
		Status:     strings.ToUpper(strings.TrimSpace(input.Status)),
		CreatedAt:  s.clock.Now(),
		UpdatedAt:  s.clock.Now(),
	}

	if problem.Status == "" {
		problem.Status = "DRAFT"
	}
	if problem.Type == "" {
		problem.Type = "SINGLE_CORRECT"
	}
	if problem.Difficulty == "" {
		problem.Difficulty = "MEDIUM"
	}
	if problem.SourceType == "" {
		problem.SourceType = "MANUAL"
	}

	if err := s.problemRepo.Save(problem); err != nil {
		return in.CreateProblemOutput{}, err
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
	if err := s.problemRepo.SaveOptions(problem.ID, options); err != nil {
		return in.CreateProblemOutput{}, err
	}

	tags := make([]domain.ProblemTag, 0, len(input.Tags))
	for _, tag := range input.Tags {
		tagName := strings.TrimSpace(tag)
		if tagName == "" {
			continue
		}
		tags = append(tags, domain.ProblemTag{ProblemID: problem.ID, Tag: tagName})
	}
	if err := s.problemRepo.SaveTags(problem.ID, tags); err != nil {
		return in.CreateProblemOutput{}, err
	}

	return in.CreateProblemOutput{ProblemID: problem.ID}, nil
}
