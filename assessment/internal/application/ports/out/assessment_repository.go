package out

import "github.com/aakashloyar/elevate/assessment/internal/domain"

type AssessmentRepository interface {
	Save(assessment domain.Assessment) error
	FindAll(filters FindAllAssessmentFilters) ([]domain.Assessment, error)
	FindByID(assessmentID string) (domain.Assessment, error)
	DeleteByID(assessmentID string) error
	AddProblems(assessmentID string, problemIDs []string) error
	FindProblemIDs(assessmentID string) ([]string, error)
	FindMarkingScheme(assessmentID string) (domain.AssessmentMarkingScheme, error)
	CreateMarkingScheme(markingScheme domain.AssessmentMarkingScheme) error
	UpsertMarkingScheme(markingScheme domain.AssessmentMarkingScheme) error
}

type FindAllAssessmentFilters struct {
	UserID      *string
	Title       *string
	Description *string
	Limit       *int
	Offset      *int
}
