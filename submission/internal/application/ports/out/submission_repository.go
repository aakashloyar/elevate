package out

import (
	"time"

	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type SubmissionRepository interface {
	Save(submission domain.Submission) error
	SaveAnswer(answer domain.SubmissionAnswer) error
	FindByID(submissionID string) (domain.Submission, []domain.SubmissionAnswer, error)
	UpdateStatus(submissionID string, status domain.SubmissionStatus) error
	UpdateStartTime(submissionID string, startedAt time.Time, status domain.SubmissionStatus) error
	UpdateSubmissionTime(submissionID string, submittedAt time.Time, status domain.SubmissionStatus) error
}
