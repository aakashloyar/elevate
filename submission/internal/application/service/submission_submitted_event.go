package submission

import (
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

func newSubmissionSubmittedMessage(topic string, submission domain.Submission, answers []domain.SubmissionAnswer) out.SubmissionSubmittedMessage {
	eventAnswers := make([]out.SubmissionAnswerEvent, 0, len(answers))
	for _, answer := range answers {
		eventAnswers = append(eventAnswers, out.SubmissionAnswerEvent{
			ProblemID: answer.ProblemID,
			Answer:    answer.Answer,
			CreatedAt: answer.CreatedAt,
			UpdatedAt: answer.UpdatedAt,
		})
	}

	return out.SubmissionSubmittedMessage{
		Topic: topic,
		Event: out.SubmissionSubmittedEvent{
			SubmissionID:    submission.ID,
			AssessmentID:    submission.AssessmentID,
			UserID:          submission.UserID,
			StartedAt:       submission.StartedAt,
			DurationSeconds: submission.DurationSeconds,
			SubmittedAt:     submission.SubmittedAt,
			Answers:         eventAnswers,
		},
	}
}
