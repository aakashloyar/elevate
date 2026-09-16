package out

import "github.com/aakashloyar/elevate/problem_generation/internal/domain"

type GenerationJobRepository interface {
	Save(job domain.GenerationJob) error
	FindByID(jobID string) (domain.GenerationJob, error)
	SaveGeneratedProblemCount(jobID string, count int) error
	UpdateStatus(jobID string, status domain.GenerationJobStatus) error
}
