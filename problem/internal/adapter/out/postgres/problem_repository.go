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

func (r *ProblemRepository) Migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS problems (
			id TEXT PRIMARY KEY,
			created_by TEXT NOT NULL,
			title TEXT NOT NULL DEFAULT '',
			statement TEXT NOT NULL,
			type TEXT NOT NULL DEFAULT 'SINGLE_CORRECT',
			difficulty TEXT NOT NULL DEFAULT 'MEDIUM',
			source_type TEXT NOT NULL DEFAULT 'MANUAL',
			status TEXT NOT NULL DEFAULT 'DRAFT',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS problem_options (
			id TEXT PRIMARY KEY,
			problem_id TEXT NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
			text TEXT NOT NULL,
			is_correct BOOLEAN NOT NULL DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS problem_tags (
			problem_id TEXT NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
			tag TEXT NOT NULL,
			PRIMARY KEY (problem_id, tag)
		)`,
	}
	for _, query := range queries {
		if _, err := r.db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProblemRepository) Save(problem domain.Problem) error {
	query := `
		INSERT INTO problems (
			id,
			created_by,
			title,
			statement,
			type,
			difficulty,
			source_type,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(query, problem.ID, problem.CreatedBy, problem.Title, problem.Statement, problem.Type, problem.Difficulty, problem.SourceType, problem.Status, problem.CreatedAt, problem.UpdatedAt)
	return err
}

func (r *ProblemRepository) SaveOptions(problemID string, options []domain.ProblemOption) error {
	if len(options) == 0 {
		return nil
	}
	stmt := `INSERT INTO problem_options (id, problem_id, text, is_correct) VALUES ($1, $2, $3, $4)`
	for _, opt := range options {
		if _, err := r.db.Exec(stmt, opt.ID, problemID, opt.Text, opt.IsCorrect); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProblemRepository) SaveTags(problemID string, tags []domain.ProblemTag) error {
	if len(tags) == 0 {
		return nil
	}
	stmt := `INSERT INTO problem_tags (problem_id, tag) VALUES ($1, $2)`
	for _, tag := range tags {
		if _, err := r.db.Exec(stmt, problemID, tag.Tag); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProblemRepository) FindByID(problemID string) (domain.Problem, []domain.ProblemOption, []domain.ProblemTag, error) {
	problemQuery := `SELECT id, created_by, title, statement, type, difficulty, source_type, status, created_at, updated_at FROM problems WHERE id = $1`
	row := r.db.QueryRow(problemQuery, problemID)

	var problem domain.Problem
	if err := row.Scan(&problem.ID, &problem.CreatedBy, &problem.Title, &problem.Statement, &problem.Type, &problem.Difficulty, &problem.SourceType, &problem.Status, &problem.CreatedAt, &problem.UpdatedAt); err != nil {
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
	query := `SELECT id, created_by, title, statement, type, difficulty, source_type, status, created_at, updated_at FROM problems`
	args := []any{}
	whereClauses := []string{}
	if value, ok := filters["type"]; ok && value != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("type = $%d", len(args)+1))
		args = append(args, value)
	}
	if value, ok := filters["status"]; ok && value != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, value)
	}
	if value, ok := filters["tag"]; ok && value != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("id IN (SELECT problem_id FROM problem_tags WHERE lower(tag) = lower($%d))", len(args)+1))
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
		if err := rows.Scan(&problem.ID, &problem.CreatedBy, &problem.Title, &problem.Statement, &problem.Type, &problem.Difficulty, &problem.SourceType, &problem.Status, &problem.CreatedAt, &problem.UpdatedAt); err != nil {
			return nil, err
		}
		problems = append(problems, problem)
	}
	return problems, nil
}

func (r *ProblemRepository) Update(problem domain.Problem) error {
	query := `
		UPDATE problems
		SET title = $2, statement = $3, type = $4, difficulty = $5, source_type = $6, status = $7, updated_at = $8
		WHERE id = $1
	`
	_, err := r.db.Exec(query, problem.ID, problem.Title, problem.Statement, problem.Type, problem.Difficulty, problem.SourceType, problem.Status, problem.UpdatedAt)
	return err
}

func (r *ProblemRepository) DeleteByID(problemID string) error {
	_, err := r.db.Exec(`DELETE FROM problems WHERE id = $1`, problemID)
	return err
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
