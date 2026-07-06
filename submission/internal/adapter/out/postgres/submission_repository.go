package postgres

import (
	"database/sql"
	"time"

	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type SubmissionRepository struct {
	db *sql.DB
}

func NewSubmissionRepository(db *sql.DB) out.SubmissionRepository {
	return &SubmissionRepository{db: db}
}

func (r *SubmissionRepository) Migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS submissions (
			id TEXT PRIMARY KEY,
			assessment_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'IN_PROGRESS',
			started_at TIMESTAMP NOT NULL,
			submitted_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS submission_answers (
			id TEXT PRIMARY KEY,
			submission_id TEXT NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
			problem_id TEXT NOT NULL,
			answer TEXT[] NOT NULL,
			answered_at TIMESTAMP NOT NULL
		)`,
	}
	for _, query := range queries {
		if _, err := r.db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func (r *SubmissionRepository) Save(submission domain.Submission) error {
	query := `
		INSERT INTO submissions (
			id,
			assessment_id,
			user_id,
			status,
			started_at,
			submitted_at,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query, submission.ID, submission.AssessmentID, submission.UserID, submission.Status, submission.StartedAt, submission.SubmittedAt, submission.CreatedAt, submission.UpdatedAt)
	return err
}

func (r *SubmissionRepository) SaveAnswer(answer domain.SubmissionAnswer) error {
	if answer.ID == "" {
		answer.ID = time.Now().UTC().Format(time.RFC3339Nano)
	}
	query := `
		INSERT INTO submission_answers (
			id,
			submission_id,
			problem_id,
			answer,
			answered_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, answer.ID, answer.SubmissionID, answer.ProblemID, pqStringArray(answer.Answer), answer.AnsweredAt)
	return err
}

func (r *SubmissionRepository) FindByID(submissionID string) (domain.Submission, []domain.SubmissionAnswer, error) {
	submissionQuery := `SELECT id, assessment_id, user_id, status, started_at, submitted_at, created_at, updated_at FROM submissions WHERE id = $1`
	row := r.db.QueryRow(submissionQuery, submissionID)

	var submission domain.Submission
	var submittedAt sql.NullTime
	if err := row.Scan(&submission.ID, &submission.AssessmentID, &submission.UserID, &submission.Status, &submission.StartedAt, &submittedAt, &submission.CreatedAt, &submission.UpdatedAt); err != nil {
		return domain.Submission{}, nil, err
	}
	if submittedAt.Valid {
		submission.SubmittedAt = &submittedAt.Time
	}

	answerQuery := `SELECT problem_id, answer, answered_at FROM submission_answers WHERE submission_id = $1 ORDER BY answered_at`
	answerRows, err := r.db.Query(answerQuery, submissionID)
	if err != nil {
		return domain.Submission{}, nil, err
	}
	defer answerRows.Close()

	answers := []domain.SubmissionAnswer{}
	for answerRows.Next() {
		var answer domain.SubmissionAnswer
		var rawAnswer []string
		if err := answerRows.Scan(&answer.ProblemID, &rawAnswer, &answer.AnsweredAt); err != nil {
			return domain.Submission{}, nil, err
		}
		answer.Answer = rawAnswer
		answers = append(answers, answer)
	}
	return submission, answers, nil
}

func (r *SubmissionRepository) UpdateStatus(submissionID, status string) error {
	_, err := r.db.Exec(`UPDATE submissions SET status = $2, updated_at = NOW() WHERE id = $1`, submissionID, status)
	return err
}

func (r *SubmissionRepository) UpdateSubmissionTime(submissionID string, submittedAt time.Time) error {
	_, err := r.db.Exec(`UPDATE submissions SET submitted_at = $2, status = 'SUBMITTED', updated_at = NOW() WHERE id = $1`, submissionID, submittedAt)
	return err
}

func pqStringArray(values []string) []string {
	return values
}
