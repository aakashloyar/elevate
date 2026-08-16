package assessment

import (
	"context"
	"errors"
	"strings"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in"
	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
)

type UpsertAssessmentMarkingSchemeService struct {
	assessmentRepo out.AssessmentRepository
}

func NewUpsertAssessmentMarkingSchemeService(assessmentRepo out.AssessmentRepository) in.UpsertAssessmentMarkingSchemeService {
	return &UpsertAssessmentMarkingSchemeService{assessmentRepo: assessmentRepo}
}

func (s *UpsertAssessmentMarkingSchemeService) Execute(ctx context.Context, input in.UpsertAssessmentMarkingSchemeInput) error {
	if strings.TrimSpace(input.AssessmentID) == "" {
		return errors.New("assessment id is required")
	}

	if err := validateMarkValues(input); err != nil {
		return err
	}

	if _, err := s.assessmentRepo.FindByID(input.AssessmentID); err != nil {
		return err
	}

	return s.assessmentRepo.UpsertMarkingScheme(toAssessmentMarkingScheme(input))
}
