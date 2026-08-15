package out

import "github.com/aakashloyar/elevate/problem/internal/domain"

type ProblemRepository interface {
	Save(problem domain.Problem, options []domain.ProblemOption, tags []domain.ProblemTag) error
	FindByID(problemID string) (domain.Problem, []domain.ProblemOption, []domain.ProblemTag, error)
	List(offset, limit int, filters map[string]string) ([]domain.Problem, error)
	Update(problem domain.Problem, options []domain.ProblemOption, tags []domain.ProblemTag) error
	DeleteByID(problemID string) error
}
