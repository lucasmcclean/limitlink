package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/lucasmcclean/url-shortener/config"
)

type Postgres struct {
	db *sql.DB
}

func New(cfg *config.DB) (*Postgres, error) {
	connStr := generateConnStr(cfg)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening databse connection: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %s", err)
	}

	return &Postgres{db: db}, nil
}

func generateConnStr(cfg *config.DB) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)
}

func (p *Postgres) Close() error {
	err := p.Close()
	return err
}
