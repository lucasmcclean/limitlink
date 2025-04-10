package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/lucasmcclean/url-shortener/config"
)

type DB struct {
	pool *sql.DB
}

func New(cfg *config.DB) (*DB, error) {
	connStr := cfg.GenerateConnStrNoSSL()

	pool, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening databse connection: %s", err)
	}

	err = pool.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging database: %s", err)
	}

	return &DB{pool: pool}, nil
}

func (db *DB) Close() error {
	err := db.pool.Close()
	return err
}
