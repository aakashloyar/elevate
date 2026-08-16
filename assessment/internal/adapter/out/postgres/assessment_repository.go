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

func (r *AssessmentRepository) FindProblemIDs(assessmentID string) ([]string, error) {
	rows, err := r.db.Query(`SELECT problem_id FROM assessment_problems WHERE assessment_id = $1 ORDER BY problem_id`, assessmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	problemIDs := make([]string, 0)
	for rows.Next() {
		var problemID string
		if err := rows.Scan(&problemID); err != nil {
			return nil, err
		}
		problemIDs = append(problemIDs, problemID)
	}
	return problemIDs, rows.Err()
}

func (r *AssessmentRepository) FindMarkingScheme(assessmentID string) (domain.AssessmentMarkingScheme, error) {
	row := r.db.QueryRow(`
		SELECT
			assessment_id,
			single_correct_marks,
			single_incorrect_marks,
			single_skipped_marks,
			multiple_correct_marks,
			multiple_incorrect_marks,
			multiple_skipped_marks,
			numerical_correct_marks,
			numerical_incorrect_marks,
			numerical_skipped_marks
		FROM assessment_marking_schemes
		WHERE assessment_id = $1`, assessmentID)

	var markingScheme domain.AssessmentMarkingScheme
	err := row.Scan(
		&markingScheme.AssessmentID,
		&markingScheme.SingleCorrectMarks,
		&markingScheme.SingleIncorrectMarks,
		&markingScheme.SingleSkippedMarks,
		&markingScheme.MultipleCorrectMarks,
		&markingScheme.MultipleIncorrectMarks,
		&markingScheme.MultipleSkippedMarks,
		&markingScheme.NumericalCorrectMarks,
		&markingScheme.NumericalIncorrectMarks,
		&markingScheme.NumericalSkippedMarks,
	)
	return markingScheme, err
}

func (r *AssessmentRepository) UpsertMarkingScheme(markingScheme domain.AssessmentMarkingScheme) error {
	_, err := r.db.Exec(`
		INSERT INTO assessment_marking_schemes (
			assessment_id,
			single_correct_marks,
			single_incorrect_marks,
			single_skipped_marks,
			multiple_correct_marks,
			multiple_incorrect_marks,
			multiple_skipped_marks,
			numerical_correct_marks,
			numerical_incorrect_marks,
			numerical_skipped_marks
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (assessment_id) DO UPDATE SET
			single_correct_marks = EXCLUDED.single_correct_marks,
			single_incorrect_marks = EXCLUDED.single_incorrect_marks,
			single_skipped_marks = EXCLUDED.single_skipped_marks,
			multiple_correct_marks = EXCLUDED.multiple_correct_marks,
			multiple_incorrect_marks = EXCLUDED.multiple_incorrect_marks,
			multiple_skipped_marks = EXCLUDED.multiple_skipped_marks,
			numerical_correct_marks = EXCLUDED.numerical_correct_marks,
			numerical_incorrect_marks = EXCLUDED.numerical_incorrect_marks,
			numerical_skipped_marks = EXCLUDED.numerical_skipped_marks`,
		markingScheme.AssessmentID,
		markingScheme.SingleCorrectMarks,
		markingScheme.SingleIncorrectMarks,
		markingScheme.SingleSkippedMarks,
		markingScheme.MultipleCorrectMarks,
		markingScheme.MultipleIncorrectMarks,
		markingScheme.MultipleSkippedMarks,
		markingScheme.NumericalCorrectMarks,
		markingScheme.NumericalIncorrectMarks,
		markingScheme.NumericalSkippedMarks,
	)
	return err
}

func (r *AssessmentRepository) CreateMarkingScheme(markingScheme domain.AssessmentMarkingScheme) error {
	_, err := r.db.Exec(`
		INSERT INTO assessment_marking_schemes (
			assessment_id,
			single_correct_marks,
			single_incorrect_marks,
			single_skipped_marks,
			multiple_correct_marks,
			multiple_incorrect_marks,
			multiple_skipped_marks,
			numerical_correct_marks,
			numerical_incorrect_marks,
			numerical_skipped_marks
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		markingScheme.AssessmentID,
		markingScheme.SingleCorrectMarks,
		markingScheme.SingleIncorrectMarks,
		markingScheme.SingleSkippedMarks,
		markingScheme.MultipleCorrectMarks,
		markingScheme.MultipleIncorrectMarks,
		markingScheme.MultipleSkippedMarks,
		markingScheme.NumericalCorrectMarks,
		markingScheme.NumericalIncorrectMarks,
		markingScheme.NumericalSkippedMarks,
	)
	return err
}
