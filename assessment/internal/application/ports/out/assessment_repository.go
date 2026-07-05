package out

import "github.com/aakashloyar/elevate/assessment/internal/domain"

type AssessmentRepository interface {
	Save(assessment domain.Assessment) error
	FindByID(assessmentID string) (domain.Assessment, error)
	DeleteByID(assessmentID string) error
}
