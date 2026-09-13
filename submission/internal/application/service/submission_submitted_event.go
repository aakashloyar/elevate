package submission

import (
	"time"

	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

func newSubmissionSubmittedMessage(topic string, submission domain.Submission, drafts []domain.SubmissionAnswerDraft) out.SubmissionSubmittedMessage {
	var startedAt time.Time
	if submission.StartedAt != nil {
		startedAt = *submission.StartedAt
	}

	eventAnswers := make([]out.SubmissionAnswerEvent, 0, len(drafts))
	for _, draft := range drafts {
		updatedAt := time.Time{}
		if draft.AnswerUpdatedAt != nil {
			updatedAt = *draft.AnswerUpdatedAt
		}
		eventAnswers = append(eventAnswers, out.SubmissionAnswerEvent{
			ProblemID: draft.ProblemID,
			Answer:    draft.Answer,
			CreatedAt: updatedAt,
			UpdatedAt: updatedAt,
		})
	}

	return out.SubmissionSubmittedMessage{
		Topic: topic,
		Event: out.SubmissionSubmittedEvent{
			SubmissionID:    submission.ID,
			AssessmentID:    submission.AssessmentID,
			UserID:          submission.UserID,
			StartedAt:       startedAt,
			DurationSeconds: submission.DurationSeconds,
			SubmittedAt:     submission.SubmittedAt,
			Answers:         eventAnswers,
		},
	}
}
