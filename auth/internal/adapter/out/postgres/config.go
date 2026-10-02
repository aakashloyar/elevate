package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Open(host, port, user, password, db, ssl string) (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, db, ssl))
}
