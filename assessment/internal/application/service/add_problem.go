package assessment

import (
	"context"
	"errors"
	"strings"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in"
	outports "github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
)

type AddProblemService struct {
	assessmentRepo outports.AssessmentRepository
	problemClient  outports.ProblemClient
}

func NewAddProblemService(assessmentRepo outports.AssessmentRepository, problemClient outports.ProblemClient) in.AddProblemService {
	return &AddProblemService{assessmentRepo: assessmentRepo, problemClient: problemClient}
}

func (s *AddProblemService) Execute(ctx context.Context, input in.AddProblemInput) (in.AddProblemOutput, error) {
	assessmentID := strings.TrimSpace(input.AssessmentID)
	if assessmentID == "" {
		return in.AddProblemOutput{}, errors.New("assessment id is required")
	}

	if _, err := s.assessmentRepo.FindByID(assessmentID); err != nil {
		return in.AddProblemOutput{}, err
	}

	options := make([]outports.CreateProblemOptionInput, 0, len(input.Options))
	for _, option := range input.Options {
		options = append(options, outports.CreateProblemOptionInput{Text: option.Text, IsCorrect: option.IsCorrect})
	}

	problemOut, err := s.problemClient.CreateProblem(ctx, outports.CreateProblemInput{
		CreatedBy:  input.CreatedBy,
		Title:      input.Title,
		Statement:  input.Statement,
		Type:       input.Type,
		Difficulty: input.Difficulty,
		SourceType: input.SourceType,
		Options:    options,
		Tags:       input.Tags,
	})
	if err != nil {
		return in.AddProblemOutput{}, err
	}

	if err := s.assessmentRepo.AddProblems(assessmentID, []string{problemOut.ProblemID}); err != nil {
		return in.AddProblemOutput{}, err
	}

	return in.AddProblemOutput{ProblemID: problemOut.ProblemID}, nil
}
