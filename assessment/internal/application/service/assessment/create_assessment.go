package assessment

import (
	"context"
	"errors"
	"strings"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in/assessment"
	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
	"github.com/aakashloyar/elevate/assessment/internal/domain"
)

type CreateAssessmentService struct {
	assessmentRepo out.AssessmentRepository
	idGen          out.IDGenerator
	clock          out.Clock
}

func NewCreateAssessmentService(assessmentRepo out.AssessmentRepository, idGen out.IDGenerator, clock out.Clock) in.CreateAssessmentService {
	return &CreateAssessmentService{assessmentRepo: assessmentRepo, idGen: idGen, clock: clock}
}

func (s *CreateAssessmentService) Execute(ctx context.Context, input in.CreateAssessmentInput) (in.CreateAssessmentOutput, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return in.CreateAssessmentOutput{}, errors.New("title is required")
	}

	description := strings.TrimSpace(input.Description)
	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "DRAFT"
	}

	if input.DurationSeconds <= 0 {
		return in.CreateAssessmentOutput{}, errors.New("duration seconds must be greater than zero")
	}

	if strings.TrimSpace(input.CreatedBy) == "" {
		return in.CreateAssessmentOutput{}, errors.New("created by is required")
	}

	now := s.clock.Now()
	assessment := domain.Assessment{
		ID:              s.idGen.NewID(),
		Title:           title,
		Description:     description,
		Status:          status,
		DurationSeconds: input.DurationSeconds,
		CreatedBy:       input.CreatedBy,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.assessmentRepo.Save(assessment); err != nil {
		return in.CreateAssessmentOutput{}, err
	}

	return in.CreateAssessmentOutput{AssessmentID: assessment.ID}, nil
}
