package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func New(getenv func(string) string) (*sql.DB, error) {
	var (
		DB_USER     = getenv("DB_USER")
		DB_PASSWORD = getenv("DB_PASSWORD")
		DB_NAME     = getenv("DB_NAME")
		DB_HOST     = getenv("DB_HOST")
		DB_PORT     = getenv("DB_PORT")
	)

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening databse connection: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %s", err)
	}

	return db, nil
}
