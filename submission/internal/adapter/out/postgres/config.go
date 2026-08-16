package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (c Config) NewDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
		c.SSLMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Neon pooler endpoints use PgBouncer transaction pooling. lib/pq sends
	// parameterized queries through the extended protocol with an unnamed
	// prepared statement. Keeping one client connection prevents concurrent
	// requests from binding against each other's unnamed statement state.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	return db, nil
}
