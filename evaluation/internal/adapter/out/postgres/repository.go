package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	evaluationservice "github.com/aakashloyar/elevate/evaluation/internal/application/service"
	"github.com/aakashloyar/elevate/evaluation/internal/domain"
	_ "github.com/lib/pq"
)

type Repository struct{ db *sql.DB }

func New(db *sql.DB) *Repository { return &Repository{db: db} }

func Open(host, port, user, password, database, sslmode string) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, database, sslmode)
	return sql.Open("postgres", dsn)
}

func (r *Repository) Save(ctx context.Context, evaluation domain.Evaluation) error {
	questions, err := json.Marshal(evaluation.Questions)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `INSERT INTO evaluations (submission_id, assessment_id, user_id, started_at, duration_seconds, submitted_at, score, questions, evaluated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (submission_id) DO NOTHING`, evaluation.SubmissionID, evaluation.AssessmentID, evaluation.UserID, evaluation.StartedAt, evaluation.DurationSeconds, evaluation.SubmittedAt, evaluation.Score, questions, evaluation.EvaluatedAt)
	return err
}
func (r *Repository) FindBySubmissionID(ctx context.Context, submissionID string) (domain.Evaluation, error) {
	var result domain.Evaluation
	var questions []byte
	err := r.db.QueryRowContext(ctx, `SELECT submission_id, assessment_id, user_id, started_at, duration_seconds, submitted_at, score, questions, evaluated_at FROM evaluations WHERE submission_id = $1`, submissionID).Scan(&result.SubmissionID, &result.AssessmentID, &result.UserID, &result.StartedAt, &result.DurationSeconds, &result.SubmittedAt, &result.Score, &questions, &result.EvaluatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Evaluation{}, evaluationservice.ErrEvaluationNotFound
	}
	if err != nil {
		return domain.Evaluation{}, err
	}
	if err := json.Unmarshal(questions, &result.Questions); err != nil {
		return domain.Evaluation{}, err
	}
	return result, nil
}
