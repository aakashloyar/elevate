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
		input.Single.Correct, input.Single.Incorrect, input.Single.Skipped,
		input.Multiple.Correct, input.Multiple.Incorrect, input.Multiple.Skipped,
		input.Numerical.Correct, input.Numerical.Incorrect, input.Numerical.Skipped,
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
		AssessmentID: input.AssessmentID,
		Single: domain.Marks{
			Correct:   input.Single.Correct,
			Incorrect: input.Single.Incorrect,
			Skipped:   input.Single.Skipped,
		},
		Multiple: domain.Marks{
			Correct:   input.Multiple.Correct,
			Incorrect: input.Multiple.Incorrect,
			Skipped:   input.Multiple.Skipped,
		},
		Numerical: domain.Marks{
			Correct:   input.Numerical.Correct,
			Incorrect: input.Numerical.Incorrect,
			Skipped:   input.Numerical.Skipped,
		},
	}
}
