package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

var ErrEvaluationNotFound = errors.New("evaluation not found")

type Service struct {
	assessments AssessmentClient
	problems    ProblemClient
	submissions SubmissionClient
	repository  Repository
	now         func() time.Time
}

func NewService(assessments AssessmentClient, problems ProblemClient, submissions SubmissionClient, repository Repository) *Service {
	return &Service{assessments: assessments, problems: problems, submissions: submissions, repository: repository, now: func() time.Time { return time.Now().UTC() }}
}

func (s *Service) Evaluate(ctx context.Context, submission domain.SubmissionSubmitted) (domain.Evaluation, error) {
	if submission.SubmissionID == "" || submission.AssessmentID == "" {
		return domain.Evaluation{}, errors.New("submission_id and assessment_id are required")
	}
	if existing, err := s.repository.FindBySubmissionID(ctx, submission.SubmissionID); err == nil {
		if err := s.submissions.UpdateStatus(ctx, submission.SubmissionID, domain.SubmissionStatusEvaluated); err != nil {
			return domain.Evaluation{}, fmt.Errorf("mark submission evaluated: %w", err)
		}
		return existing, nil
	} else if !errors.Is(err, ErrEvaluationNotFound) {
		return domain.Evaluation{}, err
	}

	if err := s.submissions.UpdateStatus(ctx, submission.SubmissionID, domain.SubmissionStatusUnderEvaluation); err != nil {
		return domain.Evaluation{}, fmt.Errorf("mark submission under evaluation: %w", err)
	}

	evaluation, err := s.evaluateSubmission(ctx, submission)
	if err != nil {
		if statusErr := s.submissions.UpdateStatus(ctx, submission.SubmissionID, domain.SubmissionStatusEvaluationFailed); statusErr != nil {
			return domain.Evaluation{}, fmt.Errorf("%w; additionally failed to mark submission evaluation failed: %v", err, statusErr)
		}
		return domain.Evaluation{}, err
	}
	if err := s.submissions.UpdateStatus(ctx, submission.SubmissionID, domain.SubmissionStatusEvaluated); err != nil {
		return domain.Evaluation{}, fmt.Errorf("mark submission evaluated: %w", err)
	}
	return evaluation, nil
}

func (s *Service) evaluateSubmission(ctx context.Context, submission domain.SubmissionSubmitted) (domain.Evaluation, error) {
	scheme, err := s.assessments.MarkingScheme(ctx, submission.AssessmentID)
	if err != nil {
		return domain.Evaluation{}, fmt.Errorf("get marking scheme: %w", err)
	}
	problemIDs, err := s.assessments.ProblemIDs(ctx, submission.AssessmentID)
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
		problem, err := s.problems.Problem(ctx, problemID)
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

func (s *Service) Get(ctx context.Context, submissionID string) (domain.Evaluation, error) {
	return s.repository.FindBySubmissionID(ctx, submissionID)
}

func EvaluateAnswer(problem domain.Problem, answer []string, scheme domain.MarkingScheme) domain.QuestionResult {
	result := domain.QuestionResult{
		ProblemID:       problem.ID,
		Type:            problem.Type,
		SelectedOptions: selectedOptions(problem.Type, problem.Options, answer),
	}
	correct, incorrect, skipped := marks(problem.Type, scheme)
	if len(answer) == 0 {
		result.Status, result.Marks = "skipped", skipped
		return result
	}
	if problem.Type == domain.ProblemTypeMultiple {
		return evaluateMultipleAnswer(result, answer, correctOptionIDs(problem.Options), correct, incorrect)
	}
	if sameSet(answer, correctOptionIDs(problem.Options)) {
		result.Status, result.Marks = "correct", correct
	} else {
		result.Status, result.Marks = "incorrect", incorrect
	}
	return result
}

func selectedOptions(problemType domain.ProblemType, options []domain.Option, answer []string) []domain.SelectedOption {
	if len(answer) == 0 {
		return []domain.SelectedOption{}
	}

	if problemType == domain.ProblemTypeNumerical {
		selected := make([]domain.SelectedOption, 0, len(answer))
		for _, value := range answer {
			selected = append(selected, domain.SelectedOption{Text: value})
		}
		return selected
	}

	optionByID := make(map[string]domain.Option, len(options))
	for _, option := range options {
		optionByID[option.ID] = option
	}

	selected := make([]domain.SelectedOption, 0, len(answer))
	for _, optionID := range answer {
		option, ok := optionByID[optionID]
		if !ok {
			continue
		}
		selected = append(selected, domain.SelectedOption{
			ID:        option.ID,
			Text:      option.Text,
			IsCorrect: option.IsCorrect,
		})
	}
	return selected
}

// evaluateMultipleAnswer gives proportional marks only when every selected
// option is correct. Selecting any incorrect option receives the configured
// incorrect mark for the question.
func evaluateMultipleAnswer(result domain.QuestionResult, answer, correctOptionIDs []string, correctMarks, incorrectMarks float64) domain.QuestionResult {
	correctOptions := make(map[string]struct{}, len(correctOptionIDs))
	for _, optionID := range correctOptionIDs {
		correctOptions[optionID] = struct{}{}
	}
	if len(correctOptions) == 0 {
		result.Status, result.Marks = "incorrect", incorrectMarks
		return result
	}

	selected := make(map[string]struct{}, len(answer))
	for _, optionID := range answer {
		if optionID == "" {
			result.Status, result.Marks = "incorrect", incorrectMarks
			return result
		}
		if _, alreadySelected := selected[optionID]; alreadySelected {
			result.Status, result.Marks = "incorrect", incorrectMarks
			return result
		}
		selected[optionID] = struct{}{}
		if _, isCorrect := correctOptions[optionID]; !isCorrect {
			result.Status, result.Marks = "incorrect", incorrectMarks
			return result
		}
	}

	if len(selected) == len(correctOptions) {
		result.Status, result.Marks = "correct", correctMarks
		return result
	}
	result.Status = "partially_correct"
	result.Marks = float64(len(selected)) * correctMarks / float64(len(correctOptions))
	return result
}

func marks(kind domain.ProblemType, s domain.MarkingScheme) (float64, float64, float64) {
	switch kind {
	case domain.ProblemTypeSingle:
		return s.SingleCorrectMarks, s.SingleIncorrectMarks, s.SingleSkippedMarks
	case domain.ProblemTypeMultiple:
		return s.MultipleCorrectMarks, s.MultipleIncorrectMarks, s.MultipleSkippedMarks
	case domain.ProblemTypeNumerical:
		return s.NumericalCorrectMarks, s.NumericalIncorrectMarks, s.NumericalSkippedMarks
	default:
		return 0, 0, 0
	}
}
func correctOptionIDs(options []domain.Option) []string {
	ids := make([]string, 0)
	for _, option := range options {
		if option.IsCorrect {
			ids = append(ids, option.ID)
		}
	}
	return ids
}
func sameSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	values := make(map[string]struct{}, len(a))
	for _, v := range a {
		if v == "" {
			return false
		}
		values[v] = struct{}{}
	}
	if len(values) != len(a) {
		return false
	}
	for _, v := range b {
		if _, ok := values[v]; !ok {
			return false
		}
	}
	return true
}
