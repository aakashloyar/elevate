package out

import (
	"time"

	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type SubmissionRepository interface {
	Save(submission domain.Submission) error
	SaveAnswer(answer domain.SubmissionAnswer) (bool, error)
	FindByID(submissionID string) (domain.Submission, []domain.SubmissionAnswer, error)
	FindStatus(submissionID string) (domain.SubmissionStatus, *time.Time, error)
	UpdateStatus(submissionID string, status domain.SubmissionStatus) error
	UpdateStartTime(submissionID string, startedAt, expiresAt time.Time, status domain.SubmissionStatus) error
	Submit(submissionID string, submittedAt time.Time) (bool, error)
	ExpireSubmissions(expiredAt time.Time, limit int) ([]string, error)
}
