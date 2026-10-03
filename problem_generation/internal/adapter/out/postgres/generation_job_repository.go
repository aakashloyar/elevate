package postgres

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/lib/pq"

	"github.com/aakashloyar/elevate/problem_generation/internal/domain"
)

type GenerationJobRepository struct {
	db *sql.DB
}

func NewGenerationJobRepository(db *sql.DB) *GenerationJobRepository {
	return &GenerationJobRepository{db: db}
}

func (r *GenerationJobRepository) Save(job domain.GenerationJob) error {
	defer observeDB("generation_jobs.save")()
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
			topic_ids,
			generated_problem_count,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		`, job.ID, job.UserID, job.SingleCorrectCount, job.MultiCorrectCount, job.NumericalCount, job.DocumentID, job.AssessmentID, job.Level, job.Description, job.Status, pqStringArray(job.TopicIDs), 0, job.CreatedAt, job.UpdatedAt)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *GenerationJobRepository) FindByID(jobID string) (domain.GenerationJob, error) {
	defer observeDB("generation_jobs.find_by_id")()
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
			topic_ids,
			generated_problem_count,
			created_at,
			updated_at
		FROM generation_jobs
		WHERE id = $1
	`

	row := r.db.QueryRow(query, jobID)

	var job domain.GenerationJob
	if err := row.Scan(&job.ID, &job.UserID, &job.SingleCorrectCount, &job.MultiCorrectCount, &job.NumericalCount, &job.DocumentID, &job.AssessmentID, &job.Level, &job.Description, &job.Status, pq.Array(&job.TopicIDs), &job.GeneratedProblemCount, &job.CreatedAt, &job.UpdatedAt); err != nil {
		return domain.GenerationJob{}, err
	}

	return job, nil
}

func (r *GenerationJobRepository) FindAll(limit, offset int, search string) ([]domain.GenerationJob, error) {
	defer observeDB("generation_jobs.find_all")()
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
			topic_ids,
			generated_problem_count,
			created_at,
			updated_at
		FROM generation_jobs
	`
	args := []any{limit, offset}
	if strings.TrimSpace(search) != "" {
		query += " WHERE description ILIKE $3"
		args = append(args, "%"+strings.TrimSpace(search)+"%")
	}
	query += " ORDER BY created_at DESC LIMIT $1 OFFSET $2"
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make([]domain.GenerationJob, 0, limit)
	for rows.Next() {
		var job domain.GenerationJob
		if err := rows.Scan(&job.ID, &job.UserID, &job.SingleCorrectCount, &job.MultiCorrectCount, &job.NumericalCount, &job.DocumentID, &job.AssessmentID, &job.Level, &job.Description, &job.Status, pq.Array(&job.TopicIDs), &job.GeneratedProblemCount, &job.CreatedAt, &job.UpdatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *GenerationJobRepository) SaveGeneratedProblemCount(jobID string, count int) error {
	defer observeDB("generation_jobs.save_problem_count")()
	result, err := r.db.Exec(`UPDATE generation_jobs SET generated_problem_count = $2, updated_at = NOW() WHERE id = $1`, jobID, count)
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

func (r *GenerationJobRepository) UpdateStatus(jobID string, status domain.GenerationJobStatus) error {
	defer observeDB("generation_jobs.update_status")()
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

func pqStringArray(values []string) []string {
	return values
}
