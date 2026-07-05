package postgres

import (
	"database/sql"

	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
	"github.com/aakashloyar/elevate/assessment/internal/domain"
)

type AssessmentRepository struct {
	db *sql.DB
}

func NewAssessmentRepository(db *sql.DB) out.AssessmentRepository {
	return &AssessmentRepository{db: db}
}

func (r *AssessmentRepository) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS assessments (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'DRAFT',
			duration_seconds INTEGER NOT NULL,
			created_by TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *AssessmentRepository) Save(assessment domain.Assessment) error {
	query := `
		INSERT INTO assessments (
			id,
			title,
			description,
			status,
			duration_seconds,
			created_by,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(query, assessment.ID, assessment.Title, assessment.Description, assessment.Status, assessment.DurationSeconds, assessment.CreatedBy, assessment.CreatedAt, assessment.UpdatedAt)
	return err
}

func (r *AssessmentRepository) FindByID(assessmentID string) (domain.Assessment, error) {
	query := `
		SELECT
			id,
			title,
			description,
			status,
			duration_seconds,
			created_by,
			created_at,
			updated_at
		FROM assessments
		WHERE id = $1
	`

	row := r.db.QueryRow(query, assessmentID)

	var assessment domain.Assessment
	if err := row.Scan(&assessment.ID, &assessment.Title, &assessment.Description, &assessment.Status, &assessment.DurationSeconds, &assessment.CreatedBy, &assessment.CreatedAt, &assessment.UpdatedAt); err != nil {
		return domain.Assessment{}, err
	}

	return assessment, nil
}

func (r *AssessmentRepository) DeleteByID(assessmentID string) error {
	_, err := r.db.Exec(`DELETE FROM assessments WHERE id = $1`, assessmentID)
	return err
}
