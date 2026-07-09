package postgres

import (
	"database/sql"
	"errors"

	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
	"github.com/aakashloyar/elevate/problem_generation/internal/domain"
)

type GenerationJobRepository struct {
	db *sql.DB
}

func NewGenerationJobRepository(db *sql.DB) out.GenerationJobRepository {
	return &GenerationJobRepository{db: db}
}

func (r *GenerationJobRepository) Migrate() error {
	statements := []string{
		`
		CREATE TABLE IF NOT EXISTS generation_jobs (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			single_correct_count INTEGER NOT NULL,
			multi_correct_count INTEGER NOT NULL,
			numerical_count INTEGER NOT NULL,
			document_id TEXT NULL,
			assessment_id TEXT NULL,
			level TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS topics (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS generation_job_topics (
			generation_job_id TEXT NOT NULL REFERENCES generation_jobs(id) ON DELETE CASCADE,
			topic_id TEXT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
			PRIMARY KEY (generation_job_id, topic_id)
		)
		`,
	}

	for _, statement := range statements {
		if _, err := r.db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}

func (r *GenerationJobRepository) Save(job domain.GenerationJob) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
		INSERT INTO generation_jobs (
			id,
			user_id,
			single_correct_count,
			multi_correct_count,
			numerical_count,
			document_id,
			assessment_id,
			level,
			description,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, job.ID, job.UserID, job.SingleCorrectCount, job.MultiCorrectCount, job.NumericalCount, job.DocumentID, job.AssessmentID, job.Level, job.Description, job.Status, job.CreatedAt, job.UpdatedAt)
	if err != nil {
		return err
	}

	if len(job.TopicIDs) > 0 {
		stmt, stmtErr := tx.Prepare(`
			INSERT INTO generation_job_topics (generation_job_id, topic_id)
			VALUES ($1, $2)
		`)
		if stmtErr != nil {
			return stmtErr
		}
		defer stmt.Close()

		for _, topicID := range job.TopicIDs {
			if _, err = stmt.Exec(job.ID, topicID); err != nil {
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *GenerationJobRepository) FindByID(jobID string) (domain.GenerationJob, error) {
	query := `
		SELECT
			id,
			user_id,
			single_correct_count,
			multi_correct_count,
			numerical_count,
			document_id,
			assessment_id,
			level,
			description,
			status,
			created_at,
			updated_at
		FROM generation_jobs
		WHERE id = $1
	`

	row := r.db.QueryRow(query, jobID)

	var job domain.GenerationJob
	if err := row.Scan(&job.ID, &job.UserID, &job.SingleCorrectCount, &job.MultiCorrectCount, &job.NumericalCount, &job.DocumentID, &job.AssessmentID, &job.Level, &job.Description, &job.Status, &job.CreatedAt, &job.UpdatedAt); err != nil {
		return domain.GenerationJob{}, err
	}

	topicRows, err := r.db.Query(`
		SELECT topic_id
		FROM generation_job_topics
		WHERE generation_job_id = $1
		ORDER BY topic_id
	`, jobID)
	if err != nil {
		return domain.GenerationJob{}, err
	}
	defer topicRows.Close()

	for topicRows.Next() {
		var topicID string
		if err := topicRows.Scan(&topicID); err != nil {
			return domain.GenerationJob{}, err
		}
		job.TopicIDs = append(job.TopicIDs, topicID)
	}

	if err := topicRows.Err(); err != nil {
		return domain.GenerationJob{}, err
	}

	return job, nil
}

func (r *GenerationJobRepository) UpdateStatus(jobID string, status string) error {
	result, err := r.db.Exec(`
		UPDATE generation_jobs
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`, jobID, status)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("generation job not found")
	}

	return nil
}
