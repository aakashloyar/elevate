package postgres

import (
	"database/sql"
	"time"

	"github.com/lib/pq"

	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type SubmissionRepository struct {
	db *sql.DB
}

func NewSubmissionRepository(db *sql.DB) out.SubmissionRepository {
	return &SubmissionRepository{db: db}
}

func (r *SubmissionRepository) Save(submission domain.Submission) error {
	query := `
		INSERT INTO submissions (
			id,
			assessment_id,
			user_id,
			status,
			started_at,
			duration_seconds,
			expires_at,
			submitted_at,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(query, submission.ID, submission.AssessmentID, submission.UserID, submission.Status, submission.StartedAt, submission.DurationSeconds, submission.ExpiresAt, submission.SubmittedAt, submission.CreatedAt, submission.UpdatedAt)
	return err
}

func (r *SubmissionRepository) SaveAnswer(answer domain.SubmissionAnswer) (bool, error) {
	if answer.ID == "" {
		answer.ID = time.Now().UTC().Format(time.RFC3339Nano)
	}
	tx, err := r.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	var status domain.SubmissionStatus
	if err := tx.QueryRow(`SELECT status FROM submissions WHERE id = $1 FOR UPDATE`, answer.SubmissionID).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if status != domain.SubmissionStatusInProgress {
		return false, nil
	}

	query := `
		INSERT INTO submission_answers (
			id,
			submission_id,
			problem_id,
			answer,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (submission_id, problem_id)
		DO UPDATE SET
			answer = EXCLUDED.answer,
			updated_at = EXCLUDED.updated_at
	`
	if _, err := tx.Exec(query, answer.ID, answer.SubmissionID, answer.ProblemID, pqStringArray(answer.Answer), answer.CreatedAt, answer.UpdatedAt); err != nil {
		return false, err
	}
	if _, err := tx.Exec(`UPDATE submissions SET updated_at = $2 WHERE id = $1`, answer.SubmissionID, answer.UpdatedAt); err != nil {
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func (r *SubmissionRepository) FindByID(submissionID string) (domain.Submission, []domain.SubmissionAnswer, error) {
	submissionQuery := `SELECT id, assessment_id, user_id, status, started_at, duration_seconds, expires_at, submitted_at, created_at, updated_at FROM submissions WHERE id = $1`
	row := r.db.QueryRow(submissionQuery, submissionID)

	var submission domain.Submission
	var submittedAt sql.NullTime
	var expiresAt sql.NullTime
	if err := row.Scan(&submission.ID, &submission.AssessmentID, &submission.UserID, &submission.Status, &submission.StartedAt, &submission.DurationSeconds, &expiresAt, &submittedAt, &submission.CreatedAt, &submission.UpdatedAt); err != nil {
		return domain.Submission{}, nil, err
	}
	if expiresAt.Valid {
		submission.ExpiresAt = &expiresAt.Time
	}
	if submittedAt.Valid {
		submission.SubmittedAt = &submittedAt.Time
	}

	answerQuery := `SELECT problem_id, answer, created_at, updated_at FROM submission_answers WHERE submission_id = $1 ORDER BY created_at`
	answerRows, err := r.db.Query(answerQuery, submissionID)
	if err != nil {
		return domain.Submission{}, nil, err
	}
	defer answerRows.Close()

	answers := []domain.SubmissionAnswer{}
	for answerRows.Next() {
		var answer domain.SubmissionAnswer
		var rawAnswer []string
		if err := answerRows.Scan(&answer.ProblemID, pq.Array(&rawAnswer), &answer.CreatedAt, &answer.UpdatedAt); err != nil {
			return domain.Submission{}, nil, err
		}
		answer.Answer = rawAnswer
		answers = append(answers, answer)
	}
	return submission, answers, nil
}

func (r *SubmissionRepository) UpdateStatus(submissionID string, status domain.SubmissionStatus) error {
	_, err := r.db.Exec(`UPDATE submissions SET status = $2, updated_at = NOW() WHERE id = $1`, submissionID, status)
	return err
}

func (r *SubmissionRepository) UpdateStartTime(submissionID string, startedAt, expiresAt time.Time, status domain.SubmissionStatus) error {
	_, err := r.db.Exec(`UPDATE submissions SET started_at = $2, expires_at = $3, status = $4, updated_at = NOW() WHERE id = $1 AND status = 'CREATED'`, submissionID, startedAt, expiresAt, status)
	return err
}

func (r *SubmissionRepository) Submit(submissionID string, submittedAt time.Time) (bool, error) {
	result, err := r.db.Exec(`UPDATE submissions
		SET status = 'SUBMITTED', submitted_at = $2, updated_at = $2
		WHERE id = $1 AND status = 'IN_PROGRESS' AND expires_at >= $2`, submissionID, submittedAt)
	if err != nil {
		return false, err
	}
	updated, err := result.RowsAffected()
	return updated == 1, err
}

func (r *SubmissionRepository) ExpireSubmissions(expiredAt time.Time, limit int) (int, error) {
	result, err := r.db.Exec(`WITH expired AS (
		SELECT id
		FROM submissions
		WHERE status = 'IN_PROGRESS' AND expires_at <= $1
		ORDER BY expires_at
		FOR UPDATE SKIP LOCKED
		LIMIT $2
	)
	UPDATE submissions AS submission
	SET status = 'SUBMITTED', submitted_at = $1, updated_at = $1
	FROM expired
	WHERE submission.id = expired.id AND submission.status = 'IN_PROGRESS'`, expiredAt, limit)
	if err != nil {
		return 0, err
	}
	updated, err := result.RowsAffected()
	return int(updated), err
}

func pqStringArray(values []string) []string {
	return values
}
