package postgres

import (
	"database/sql"
	"github.com/aakashloyar/elevate/user/internal/application/ports/out"
	"github.com/aakashloyar/elevate/user/internal/domain"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) FindByIDs(userIDs []string) ([]domain.User, error) {
	defer observeDB("users.find_by_ids")()
	if len(userIDs) == 0 {
		return []domain.User{}, nil
	}
	rows, err := r.db.Query(`SELECT id, username, email, created_at, updated_at FROM users WHERE id = ANY($1) ORDER BY id`, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]domain.User, 0, len(userIDs))
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func NewUserRepository(db *sql.DB) out.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(user domain.User) error {
	defer observeDB("users.save")()
	query := `
		INSERT INTO users (
			id,
			username,
			email,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) FindByID(userID string) (domain.User, error) {
	defer observeDB("users.find_by_id")()
	query := `
		SELECT
			id,
			username,
			email,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRow(query, userID)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) Delete(userID string) (sql.Result, error) {
	defer observeDB("users.delete")()
	result, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, userID)
	return result, err
}
