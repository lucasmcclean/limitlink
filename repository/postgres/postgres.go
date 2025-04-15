package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/lucasmcclean/url-shortener/config"
	"github.com/lucasmcclean/url-shortener/repository"
)

type DB struct {
	pool    *sql.DB
	connStr string
}

func New(cfg *config.DB) (repository.Repository, error) {
	connStr := cfg.GenerateConnStrNoSSL()

	pool, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening databse connection: %s", err)
	}

	err = pool.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging database: %s", err)
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
