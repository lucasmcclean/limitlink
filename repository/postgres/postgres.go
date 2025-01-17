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
	connStr := generateConnStr(cfg)

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

func generateConnStr(cfg *config.DB) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)
}

func (db *DB) Close() error {
	err := db.pool.Close()
	return err
}
