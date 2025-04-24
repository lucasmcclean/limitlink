package config

import (
	"fmt"
	"os"

	"github.com/lucasmcclean/url-shortener/internal/logger"
)

type DB struct {
	URL string
}

func GetDB(log logger.Logger) *DB {
	dbCfg := &DB{}

	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")

	missing := []string{}

	if host == "" {
		missing = append(missing, "DB_HOST")
	}
	if port == "" {
		missing = append(missing, "DB_PORT")
	}
	if name == "" {
		missing = append(missing, "DB_NAME")
	}
	if user == "" {
		missing = append(missing, "DB_USER")
	}
	if pass == "" {
		missing = append(missing, "DB_PASS")
	}

	if len(missing) > 0 {
		log.Fatal("missing required database environment variables", "missing_env_vars", missing)
	}

	dbCfg.URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, name)

	return dbCfg
}

func (db *DB) GenerateConnStr() string {
	return db.URL
}

func (db *DB) GenerateConnStrNoSSL() string {
	return fmt.Sprintf("%s?sslmode=disable", db.URL)
}
