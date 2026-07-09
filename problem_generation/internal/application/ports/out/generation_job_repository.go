package out

import "github.com/aakashloyar/elevate/problem_generation/internal/domain"

type GenerationJobRepository interface {
	Save(job domain.GenerationJob) error
	FindByID(jobID string) (domain.GenerationJob, error)
	UpdateStatus(jobID string, status string) error
}
