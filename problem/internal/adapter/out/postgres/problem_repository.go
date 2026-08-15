package postgres

import (
	"database/sql"
	"fmt"

	"github.com/aakashloyar/elevate/problem/internal/application/ports/out"
	"github.com/aakashloyar/elevate/problem/internal/domain"
)

type ProblemRepository struct {
	db *sql.DB
}

func NewProblemRepository(db *sql.DB) out.ProblemRepository {
	return &ProblemRepository{db: db}
}

func (r *ProblemRepository) Save(problem domain.Problem, options []domain.ProblemOption, tags []domain.ProblemTag) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	problemQuery := `
		INSERT INTO problems (
			id,
			created_by,
			title,
			statement,
			type,
			difficulty,
			source_type,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	if _, err := tx.Exec(problemQuery, problem.ID, problem.CreatedBy, problem.Title, problem.Statement, problem.Type, problem.Difficulty, problem.SourceType, problem.CreatedAt, problem.UpdatedAt); err != nil {
		return err
	}

	if len(options) > 0 {
		optionStmt := `INSERT INTO problem_options (id, problem_id, text, is_correct) VALUES ($1, $2, $3, $4)`
		for _, opt := range options {
			if _, err := tx.Exec(optionStmt, opt.ID, problem.ID, opt.Text, opt.IsCorrect); err != nil {
				return err
			}
		}
	}

	if len(tags) > 0 {
		tagStmt := `INSERT INTO problem_tags (problem_id, tag) VALUES ($1, $2)`
		for _, tag := range tags {
			if _, err := tx.Exec(tagStmt, problem.ID, tag.Tag); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *ProblemRepository) FindByID(problemID string) (domain.Problem, []domain.ProblemOption, []domain.ProblemTag, error) {
	problemQuery := `SELECT id, created_by, title, statement, type, difficulty, source_type, created_at, updated_at FROM problems WHERE id = $1`
	row := r.db.QueryRow(problemQuery, problemID)

	var problem domain.Problem
	if err := row.Scan(&problem.ID, &problem.CreatedBy, &problem.Title, &problem.Statement, &problem.Type, &problem.Difficulty, &problem.SourceType, &problem.CreatedAt, &problem.UpdatedAt); err != nil {
		return domain.Problem{}, nil, nil, err
	}

	optionsQuery := `SELECT id, text, is_correct FROM problem_options WHERE problem_id = $1 ORDER BY id`
	optionRows, err := r.db.Query(optionsQuery, problemID)
	if err != nil {
		return domain.Problem{}, nil, nil, err
	}
	defer optionRows.Close()

	options := []domain.ProblemOption{}
	for optionRows.Next() {
		var opt domain.ProblemOption
		if err := optionRows.Scan(&opt.ID, &opt.Text, &opt.IsCorrect); err != nil {
			return domain.Problem{}, nil, nil, err
		}
		opt.ProblemID = problemID
		options = append(options, opt)
	}

	tagsQuery := `SELECT tag FROM problem_tags WHERE problem_id = $1 ORDER BY tag`
	tagRows, err := r.db.Query(tagsQuery, problemID)
	if err != nil {
		return domain.Problem{}, nil, nil, err
	}
	defer tagRows.Close()

	tags := []domain.ProblemTag{}
	for tagRows.Next() {
		var tag domain.ProblemTag
		if err := tagRows.Scan(&tag.Tag); err != nil {
			return domain.Problem{}, nil, nil, err
		}
		tag.ProblemID = problemID
		tags = append(tags, tag)
	}

	return problem, options, tags, nil
}

func (r *ProblemRepository) List(offset, limit int, filters map[string]string) ([]domain.Problem, error) {
	query := `SELECT id, created_by, title, statement, type, difficulty, source_type, created_at, updated_at FROM problems`
	args := []any{}
	whereClauses := []string{}
	if value, ok := filters["created_by"]; ok && value != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("created_by = $%d", len(args)+1))
		args = append(args, value)
	}
	if value, ok := filters["title"]; ok && value != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("title ILIKE $%d", len(args)+1))
		args = append(args, "%"+value+"%")
	}
	if value, ok := filters["type"]; ok && value != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("type = $%d", len(args)+1))
		args = append(args, value)
	}
	if value, ok := filters["difficulty"]; ok && value != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("difficulty = $%d", len(args)+1))
		args = append(args, value)
	}
	if value, ok := filters["source_type"]; ok && value != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("source_type = $%d", len(args)+1))
		args = append(args, value)
	}
	if value, ok := filters["tag"]; ok && value != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("EXISTS (SELECT 1 FROM problem_tags pt WHERE pt.problem_id = problems.id AND lower(pt.tag) = lower($%d))", len(args)+1))
		args = append(args, value)
	}
	if len(whereClauses) > 0 {
		query += " WHERE " + joinClauses(whereClauses)
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	problems := []domain.Problem{}
	for rows.Next() {
		var problem domain.Problem
		if err := rows.Scan(&problem.ID, &problem.CreatedBy, &problem.Title, &problem.Statement, &problem.Type, &problem.Difficulty, &problem.SourceType, &problem.CreatedAt, &problem.UpdatedAt); err != nil {
			return nil, err
		}
		problems = append(problems, problem)
	}
	return problems, nil
}

func (r *ProblemRepository) Update(problem domain.Problem, options []domain.ProblemOption, tags []domain.ProblemTag) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE problems
		SET created_by = $2, title = $3, statement = $4, type = $5, difficulty = $6, source_type = $7, updated_at = $8
		WHERE id = $1
	`
	if _, err := tx.Exec(query, problem.ID, problem.CreatedBy, problem.Title, problem.Statement, problem.Type, problem.Difficulty, problem.SourceType, problem.UpdatedAt); err != nil {
		return err
	}

	if _, err := tx.Exec(`DELETE FROM problem_options WHERE problem_id = $1`, problem.ID); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM problem_tags WHERE problem_id = $1`, problem.ID); err != nil {
		return err
	}

	if len(options) > 0 {
		optionStmt := `INSERT INTO problem_options (id, problem_id, text, is_correct) VALUES ($1, $2, $3, $4)`
		for _, opt := range options {
			if _, err := tx.Exec(optionStmt, opt.ID, problem.ID, opt.Text, opt.IsCorrect); err != nil {
				return err
			}
		}
	}

	if len(tags) > 0 {
		tagStmt := `INSERT INTO problem_tags (problem_id, tag) VALUES ($1, $2)`
		for _, tag := range tags {
			if _, err := tx.Exec(tagStmt, problem.ID, tag.Tag); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *ProblemRepository) DeleteByID(problemID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM problem_options WHERE problem_id = $1`, problemID); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM problem_tags WHERE problem_id = $1`, problemID); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM problems WHERE id = $1`, problemID); err != nil {
		return err
	}

	return tx.Commit()
}

func joinClauses(clauses []string) string {
	if len(clauses) == 0 {
		return ""
	}
	result := clauses[0]
	for _, clause := range clauses[1:] {
		result += " AND " + clause
	}
	return result
}
