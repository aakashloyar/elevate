package postgres

import (
	"database/sql"
	"github.com/aakashloyar/elevate/user/internal/application/ports/out"
	"github.com/aakashloyar/elevate/user/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) out.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(user domain.User) error {
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

func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE lower(username) = lower($1))`, username).Scan(&exists)
	return exists, err
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE lower(email) = lower($1))`, email).Scan(&exists)
	return exists, err
}
