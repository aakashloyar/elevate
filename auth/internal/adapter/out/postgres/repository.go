package postgres

import (
	"database/sql"
	"time"
)

type Credential struct {
	UserID   string
	Username string
	Email    string
}

type OTP struct {
	ID        string
	UserID    string
	CodeHash  string
	ExpiresAt time.Time
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(credential Credential) error {
	defer observeDB("auth_credentials.save")()
	query := `
		INSERT INTO auth_credentials (user_id, username, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
	`
	_, err := r.db.Exec(query, credential.UserID, credential.Username, credential.Email, time.Now().UTC())
	return err
}

func (r *Repository) FindByIdentifier(identifier string) (Credential, error) {
	defer observeDB("auth_credentials.find_by_identifier")()
	query := `
		SELECT user_id, username, email
		FROM auth_credentials
		WHERE lower(email) = lower($1) OR lower(username) = lower($1)
	`

	var credential Credential
	err := r.db.QueryRow(query, identifier).Scan(&credential.UserID, &credential.Username, &credential.Email)
	return credential, err
}

func (r *Repository) FindByUserID(userID string) (Credential, error) {
	defer observeDB("auth_credentials.find_by_user_id")()
	query := `SELECT user_id, username, email FROM auth_credentials WHERE user_id = $1`
	var credential Credential
	err := r.db.QueryRow(query, userID).Scan(&credential.UserID, &credential.Username, &credential.Email)
	return credential, err
}

func (r *Repository) SaveOTP(otp OTP) error {
	defer observeDB("email_otps.save")()
	query := `
		INSERT INTO email_otps (id, user_id, code_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, otp.ID, otp.UserID, otp.CodeHash, otp.ExpiresAt, time.Now().UTC())
	return err
}

func (r *Repository) ConsumeOTP(userID, codeHash string) (bool, error) {
	defer observeDB("email_otps.consume")()
	tx, err := r.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	var otpID string
	query := `
		SELECT id
		FROM email_otps
		WHERE user_id = $1
		  AND code_hash = $2
		  AND consumed_at IS NULL
		  AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
		FOR UPDATE
	`
	if err := tx.QueryRow(query, userID, codeHash).Scan(&otpID); err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if _, err := tx.Exec(`UPDATE email_otps SET consumed_at = $2 WHERE id = $1`, otpID, time.Now().UTC()); err != nil {
		return false, err
	}
	return true, tx.Commit()
}
