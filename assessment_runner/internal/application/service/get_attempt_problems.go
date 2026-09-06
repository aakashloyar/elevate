package service

import (
	"context"
	"fmt"

	in "github.com/aakashloyar/elevate/assessment_runner/internal/application/ports/in"
	"github.com/aakashloyar/elevate/assessment_runner/internal/application/ports/out"
)

type GetAttemptProblemsService struct {
	submissions out.SubmissionClient
	assessments out.AssessmentClient
	problems    out.ProblemClient
}

func NewGetAttemptProblemsService(submissions out.SubmissionClient, assessments out.AssessmentClient, problems out.ProblemClient) in.GetAttemptProblemsService {
	return &GetAttemptProblemsService{submissions: submissions, assessments: assessments, problems: problems}
}

func (s *GetAttemptProblemsService) Execute(ctx context.Context, input in.GetAttemptProblemsInput) (in.GetAttemptProblemsOutput, error) {
	if input.AttemptID == "" {
		return in.GetAttemptProblemsOutput{}, fmt.Errorf("attempt_id is required")
	}
	if input.Offset < 0 || input.Limit <= 0 {
		return in.GetAttemptProblemsOutput{}, fmt.Errorf("invalid pagination")
	}

	assessmentID, err := s.submissions.GetAttemptAssessmentID(ctx, input.AttemptID)
	if err != nil {
		return in.GetAttemptProblemsOutput{}, fmt.Errorf("resolve attempt: %w", err)
	}
	problemIDs, err := s.assessments.GetAssessmentProblemIDs(ctx, assessmentID)
	if err != nil {
		return in.GetAttemptProblemsOutput{}, fmt.Errorf("get assessment problems: %w", err)
	}
	if input.Offset >= len(problemIDs) {
		return in.GetAttemptProblemsOutput{Problems: []in.ProblemView{}}, nil
	}

	end := input.Offset + input.Limit
	if end > len(problemIDs) {
		end = len(problemIDs)
	}
	problems := make([]in.ProblemView, 0, end-input.Offset)
	for _, problemID := range problemIDs[input.Offset:end] {
		problem, err := s.problems.GetProblemByID(ctx, problemID)
		if err != nil {
			return in.GetAttemptProblemsOutput{}, fmt.Errorf("get problem %s: %w", problemID, err)
		}
		problems = append(problems, problem)
	}
	return in.GetAttemptProblemsOutput{Problems: problems}, nil
}
