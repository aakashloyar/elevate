package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
	"github.com/aakashloyar/elevate/assessment/internal/domain"
)

type AssessmentRepository struct {
	db *sql.DB
}

func NewAssessmentRepository(db *sql.DB) out.AssessmentRepository {
	return &AssessmentRepository{db: db}
}

func (r *AssessmentRepository) FindAll(filters out.FindAllAssessmentFilters) ([]domain.Assessment, error) {
	query := `
		SELECT
			id,
			title,
			description,
			duration_seconds,
			created_by,
			created_at,
			updated_at
		FROM assessments
		WHERE 1 = 1`
	args := make([]any, 0)
	argIndex := 1

	if filters.UserID != nil {
		query += fmt.Sprintf(" AND created_by = $%d", argIndex)
		args = append(args, *filters.UserID)
		argIndex++
	}
	if filters.Title != nil {
		query += fmt.Sprintf(" AND title ILIKE $%d", argIndex)
		args = append(args, "%"+strings.TrimSpace(*filters.Title)+"%")
		argIndex++
	}
	if filters.Description != nil {
		query += fmt.Sprintf(" AND description ILIKE $%d", argIndex)
		args = append(args, "%"+strings.TrimSpace(*filters.Description)+"%")
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	if filters.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, *filters.Limit)
		argIndex++
	}
	if filters.Offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, *filters.Offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assessments := make([]domain.Assessment, 0)
	for rows.Next() {
		var assessment domain.Assessment
		if err := rows.Scan(&assessment.ID, &assessment.Title, &assessment.Description, &assessment.DurationSeconds, &assessment.CreatedBy, &assessment.CreatedAt, &assessment.UpdatedAt); err != nil {
			return nil, err
		}
		assessments = append(assessments, assessment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assessments, nil
}

func (r *AssessmentRepository) Save(assessment domain.Assessment) error {
	query := `
		INSERT INTO assessments (
			id,
			title,
			description,
			duration_seconds,
			created_by,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(query, assessment.ID, assessment.Title, assessment.Description, assessment.DurationSeconds, assessment.CreatedBy, assessment.CreatedAt, assessment.UpdatedAt)
	return err
}

func (r *AssessmentRepository) FindByID(assessmentID string) (domain.Assessment, error) {
	query := `
		SELECT
			id,
			title,
			description,
			duration_seconds,
			created_by,
			created_at,
			updated_at
		FROM assessments
		WHERE id = $1
	`

	row := r.db.QueryRow(query, assessmentID)

	var assessment domain.Assessment
	if err := row.Scan(&assessment.ID, &assessment.Title, &assessment.Description, &assessment.DurationSeconds, &assessment.CreatedBy, &assessment.CreatedAt, &assessment.UpdatedAt); err != nil {
		return domain.Assessment{}, err
	}

	return assessment, nil
}

func (r *AssessmentRepository) DeleteByID(assessmentID string) error {
	_, err := r.db.Exec(`DELETE FROM assessments WHERE id = $1`, assessmentID)
	return err
}

func (r *AssessmentRepository) AddProblems(assessmentID string, problemIDs []string) error {
	if len(problemIDs) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO assessment_problems (assessment_id, problem_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, pid := range problemIDs {
		if _, err := stmt.Exec(assessmentID, pid); err != nil {
			return err
		}
	}

	return tx.Commit()
}
