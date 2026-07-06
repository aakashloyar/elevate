package out
package out

import "github.com/aakashloyar/elevate/submission/internal/domain"

type SubmissionRepository interface {
	Save(submission domain.Submission) error
	SaveAnswer(answer domain.SubmissionAnswer) error
	FindByID(submissionID string) (domain.Submission, []domain.SubmissionAnswer, error)
	UpdateStatus(submissionID, status string) error
	UpdateSubmissionTime(submissionID string, submittedAt time.Time) error
}
