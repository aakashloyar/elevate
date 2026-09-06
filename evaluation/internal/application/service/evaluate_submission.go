package evaluation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aakashloyar/elevate/evaluation/internal/application/ports/out"
	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

type EvaluateSubmissionService struct {
	assessments out.AssessmentClient
	problems    out.ProblemClient
	submissions out.SubmissionClient
	repository  out.Repository
	now         func() time.Time
}

func NewEvaluateSubmissionService(assessments out.AssessmentClient, problems out.ProblemClient, submissions out.SubmissionClient, repository out.Repository) *EvaluateSubmissionService {
	return &EvaluateSubmissionService{
		assessments: assessments,
		problems:    problems,
		submissions: submissions,
		repository:  repository,
		now:         func() time.Time { return time.Now().UTC() },
	}
}

func (s *EvaluateSubmissionService) Execute(ctx context.Context, submission domain.SubmissionSubmitted) (domain.Evaluation, error) {
	if submission.SubmissionID == "" || submission.AssessmentID == "" {
		return domain.Evaluation{}, errors.New("submission_id and assessment_id are required")
	}

	if existing, err := s.repository.FindBySubmissionID(ctx, submission.SubmissionID); err == nil {
		if err := s.submissions.UpdateSubmissionStatus(ctx, submission.SubmissionID, domain.SubmissionStatusEvaluated); err != nil {
			return domain.Evaluation{}, fmt.Errorf("mark submission evaluated: %w", err)
		}
		return existing, nil
	} else if !errors.Is(err, ErrEvaluationNotFound) {
		return domain.Evaluation{}, err
	}

	if err := s.submissions.UpdateSubmissionStatus(ctx, submission.SubmissionID, domain.SubmissionStatusUnderEvaluation); err != nil {
		return domain.Evaluation{}, fmt.Errorf("mark submission under evaluation: %w", err)
	}

	evaluation, err := s.evaluateSubmission(ctx, submission)
	if err != nil {
		if statusErr := s.submissions.UpdateSubmissionStatus(ctx, submission.SubmissionID, domain.SubmissionStatusEvaluationFailed); statusErr != nil {
			return domain.Evaluation{}, fmt.Errorf("%w; additionally failed to mark submission evaluation failed: %v", err, statusErr)
		}
		return domain.Evaluation{}, err
	}

	if err := s.submissions.UpdateSubmissionStatus(ctx, submission.SubmissionID, domain.SubmissionStatusEvaluated); err != nil {
		return domain.Evaluation{}, fmt.Errorf("mark submission evaluated: %w", err)
	}
	return evaluation, nil
}

func (s *EvaluateSubmissionService) evaluateSubmission(ctx context.Context, submission domain.SubmissionSubmitted) (domain.Evaluation, error) {
	scheme, err := s.assessments.GetAssessmentMarkingScheme(ctx, submission.AssessmentID)
	if err != nil {
		return domain.Evaluation{}, fmt.Errorf("get marking scheme: %w", err)
	}

	problemIDs, err := s.assessments.GetAssessmentProblemIDs(ctx, submission.AssessmentID)
	if err != nil {
		return domain.Evaluation{}, fmt.Errorf("get assessment problems: %w", err)
	}

	answers := make(map[string][]string, len(submission.Answers))
	for _, answer := range submission.Answers {
		answers[answer.ProblemID] = answer.Answer
	}

	result := domain.Evaluation{
		SubmissionID:    submission.SubmissionID,
		AssessmentID:    submission.AssessmentID,
		UserID:          submission.UserID,
		StartedAt:       submission.StartedAt,
		DurationSeconds: submission.DurationSeconds,
		SubmittedAt:     submission.SubmittedAt,
		Questions:       make([]domain.QuestionResult, 0, len(problemIDs)),
		EvaluatedAt:     s.now(),
	}

	for _, problemID := range problemIDs {
		problem, err := s.problems.GetProblemByID(ctx, problemID)
		if err != nil {
			return domain.Evaluation{}, fmt.Errorf("get problem %s: %w", problemID, err)
		}
		question := EvaluateAnswer(problem, answers[problemID], scheme)
		result.Questions = append(result.Questions, question)
		result.Score += question.Marks
	}

	if err := s.repository.Save(ctx, result); err != nil {
		return domain.Evaluation{}, fmt.Errorf("save evaluation: %w", err)
	}
	return result, nil
}
