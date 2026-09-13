package out

import (
	"time"

	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type SubmissionRepository interface {
	Save(submission domain.Submission, drafts []domain.SubmissionAnswerDraft) error
	SaveAnswer(answer domain.SubmissionAnswerDraft) (bool, error)
	SaveAnswers(answers []domain.SubmissionAnswerDraft) (bool, error)
	FindAnswerSnapshots(submissionID string) (domain.SubmissionStatus, []domain.SubmissionAnswerDraft, error)
	FindByID(submissionID string) (domain.Submission, []domain.SubmissionAnswerDraft, error)
	FindStatus(submissionID string) (domain.SubmissionStatus, *time.Time, error)
	UpdateStatus(submissionID string, status domain.SubmissionStatus) error
	UpdateStartTime(submissionID string, startedAt, expiresAt time.Time, status domain.SubmissionStatus) error
	Submit(submissionID string, submittedAt time.Time) (bool, error)
	ExpireSubmissions(expiredAt time.Time, limit int) ([]string, error)
}
