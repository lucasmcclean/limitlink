package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/lucasmcclean/url-shortener/config"
)

type DB struct {
	pool    *sql.DB
	connStr string
}

// It is reccomended to ping the database to ensure a proper connection.
func New(cfg *config.DB) (*DB, error) {
	connStr := cfg.GenerateConnStrNoSSL()

	pool, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening databse connection: %s", err)
	}

	return &DB{
		pool:    pool,
		connStr: connStr,
	}, nil
}

func (db *DB) Close() error {
	err := db.pool.Close()
	return err
}

func (db *DB) Ping() error {
	return db.pool.Ping()
}
