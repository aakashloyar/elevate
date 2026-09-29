package submission

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type CreateSubmissionService struct {
	submissionRepo   out.SubmissionRepository
	problemClient    out.ProblemClient
	assessmentClient out.AssessmentClient
	userClient       out.UserClient
	idGen            out.IDGenerator
	clock            out.Clock
}

func NewCreateSubmissionService(submissionRepo out.SubmissionRepository, problemClient out.ProblemClient, assessmentClient out.AssessmentClient, userClient out.UserClient, idGen out.IDGenerator, clock out.Clock) in.CreateSubmissionService {
	return &CreateSubmissionService{submissionRepo: submissionRepo, problemClient: problemClient, assessmentClient: assessmentClient, userClient: userClient, idGen: idGen, clock: clock}
}

func (s *CreateSubmissionService) Execute(ctx context.Context, input in.CreateSubmissionInput) (in.CreateSubmissionOutput, error) {
	if strings.TrimSpace(input.AssessmentID) == "" {
		return in.CreateSubmissionOutput{}, errors.New("assessment id is required")
	}
	if strings.TrimSpace(input.UserID) == "" {
		return in.CreateSubmissionOutput{}, errors.New("user id is required")
	}
	if input.DurationSeconds <= 0 {
		return in.CreateSubmissionOutput{}, errors.New("duration seconds must be greater than zero")
	}
	if input.MarkingScheme.AssessmentID != "" && input.MarkingScheme.AssessmentID != input.AssessmentID {
		return in.CreateSubmissionOutput{}, errors.New("marking scheme assessment id does not match assessment")
	}
	input.MarkingScheme.AssessmentID = input.AssessmentID
	if err := validateMarkingScheme(input.MarkingScheme); err != nil {
		return in.CreateSubmissionOutput{}, err
	}
	if err := s.userClient.Exists(ctx, input.UserID); err != nil {
		return in.CreateSubmissionOutput{}, fmt.Errorf("user id is invalid: %w", err)
	}
	problemIDs, err := s.assessmentClient.GetProblemIDs(ctx, input.AssessmentID)
	if err != nil {
		return in.CreateSubmissionOutput{}, fmt.Errorf("load assessment problems: %w", err)
	}
	if len(problemIDs) == 0 {
		return in.CreateSubmissionOutput{}, errors.New("assessment has no problems")
	}
	problemSnapshots, err := s.problemClient.GetProblemSnapshots(ctx, problemIDs)
	if err != nil {
		return in.CreateSubmissionOutput{}, fmt.Errorf("load assessment problems: %w", err)
	}

	now := s.clock.Now()
	submission := domain.Submission{
		ID:              s.idGen.NewID(),
		AssessmentID:    input.AssessmentID,
		UserID:          input.UserID,
		Status:          domain.SubmissionStatusCreated,
		DurationSeconds: input.DurationSeconds,
		CreatedAt:       now,
		UpdatedAt:       now,
		MarkingScheme:   &input.MarkingScheme,
	}

	drafts := make([]domain.SubmissionAnswerDraft, 0, len(problemSnapshots))
	for _, snapshot := range problemSnapshots {
		options := make([]domain.Option, 0, len(snapshot.Options))
		for _, option := range snapshot.Options {
			options = append(options, domain.Option{
				ID:   option.ID,
				Text: option.Text,
			})
		}
		drafts = append(drafts, domain.SubmissionAnswerDraft{
			SubmissionID: submission.ID,
			ProblemID:    snapshot.ProblemID,
			ProblemType:  snapshot.ProblemType,
			Options:      options,
		})
	}
	if err := s.submissionRepo.Save(submission, drafts); err != nil {
		return in.CreateSubmissionOutput{}, err
	}

	return in.CreateSubmissionOutput{SubmissionID: submission.ID, CreatedAt: now.Format(time.RFC3339)}, nil
}

func validateMarkingScheme(scheme domain.MarkingScheme) error {
	marks := []float64{
		scheme.Single.Correct, scheme.Single.Incorrect, scheme.Single.Skipped,
		scheme.Multiple.Correct, scheme.Multiple.Incorrect, scheme.Multiple.Skipped,
		scheme.Numerical.Correct, scheme.Numerical.Incorrect, scheme.Numerical.Skipped,
	}
	for _, mark := range marks {
		if math.IsNaN(mark) || math.IsInf(mark, 0) {
			return errors.New("mark values must be finite numbers")
		}
	}
	return nil
}
