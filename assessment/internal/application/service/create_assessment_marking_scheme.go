package assessment

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"strings"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in"
	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
	"github.com/aakashloyar/elevate/assessment/internal/domain"
)

var ErrAssessmentMarkingSchemeAlreadyExists = errors.New("assessment marking scheme already exists")

type CreateAssessmentMarkingSchemeService struct {
	assessmentRepo out.AssessmentRepository
}

func NewCreateAssessmentMarkingSchemeService(assessmentRepo out.AssessmentRepository) in.CreateAssessmentMarkingSchemeService {
	return &CreateAssessmentMarkingSchemeService{assessmentRepo: assessmentRepo}
}

func (s *CreateAssessmentMarkingSchemeService) Execute(ctx context.Context, input in.UpsertAssessmentMarkingSchemeInput) error {
	if strings.TrimSpace(input.AssessmentID) == "" {
		return errors.New("assessment id is required")
	}
	if err := validateMarkValues(input); err != nil {
		return err
	}
	if _, err := s.assessmentRepo.FindByID(input.AssessmentID); err != nil {
		return err
	}
	if _, err := s.assessmentRepo.FindMarkingScheme(input.AssessmentID); err == nil {
		return ErrAssessmentMarkingSchemeAlreadyExists
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	return s.assessmentRepo.CreateMarkingScheme(toAssessmentMarkingScheme(input))
}

func validateMarkValues(input in.UpsertAssessmentMarkingSchemeInput) error {
	marks := []float64{
		input.SingleCorrectMarks, input.SingleIncorrectMarks, input.SingleSkippedMarks,
		input.MultipleCorrectMarks, input.MultipleIncorrectMarks, input.MultipleSkippedMarks,
		input.NumericalCorrectMarks, input.NumericalIncorrectMarks, input.NumericalSkippedMarks,
	}
	for _, mark := range marks {
		if math.IsNaN(mark) || math.IsInf(mark, 0) {
			return errors.New("mark values must be finite numbers")
		}
	}
	return nil
}

func toAssessmentMarkingScheme(input in.UpsertAssessmentMarkingSchemeInput) domain.AssessmentMarkingScheme {
	return domain.AssessmentMarkingScheme{
		AssessmentID:            input.AssessmentID,
		SingleCorrectMarks:      input.SingleCorrectMarks,
		SingleIncorrectMarks:    input.SingleIncorrectMarks,
		SingleSkippedMarks:      input.SingleSkippedMarks,
		MultipleCorrectMarks:    input.MultipleCorrectMarks,
		MultipleIncorrectMarks:  input.MultipleIncorrectMarks,
		MultipleSkippedMarks:    input.MultipleSkippedMarks,
		NumericalCorrectMarks:   input.NumericalCorrectMarks,
		NumericalIncorrectMarks: input.NumericalIncorrectMarks,
		NumericalSkippedMarks:   input.NumericalSkippedMarks,
	}
}
