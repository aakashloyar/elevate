package problem

import (
	"context"

	in "github.com/aakashloyar/elevate/problem/internal/application/ports/in/problem"
	"github.com/aakashloyar/elevate/problem/internal/application/ports/out"
)

type GetProblemService struct {
	problemRepo out.ProblemRepository
}

func NewGetProblemService(problemRepo out.ProblemRepository) in.GetProblemService {
	return &GetProblemService{problemRepo: problemRepo}
}

func (s *GetProblemService) Execute(ctx context.Context, input in.GetProblemInput) (in.GetProblemOutput, error) {
	problem, options, tags, err := s.problemRepo.FindByID(input.ProblemID)
	if err != nil {
		return in.GetProblemOutput{}, err
	}

	optionOutputs := make([]in.ProblemOptionOutput, 0, len(options))
	for _, option := range options {
		optionOutputs = append(optionOutputs, in.ProblemOptionOutput{ID: option.ID, Text: option.Text, IsCorrect: option.IsCorrect})
	}

	tagsOutput := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagsOutput = append(tagsOutput, tag.Tag)
	}

	return in.GetProblemOutput{
		ID:         problem.ID,
		CreatedBy:  problem.CreatedBy,
		Title:      problem.Title,
		Statement:  problem.Statement,
		Type:       problem.Type,
		Difficulty: problem.Difficulty,
		SourceType: problem.SourceType,
		Status:     problem.Status,
		Options:    optionOutputs,
		Tags:       tagsOutput,
		CreatedAt:  problem.CreatedAt,
		UpdatedAt:  problem.UpdatedAt,
	}, nil
}
