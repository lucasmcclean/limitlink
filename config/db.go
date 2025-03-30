package config

import (
	"fmt"
	"os"

	"github.com/lucasmcclean/url-shortener/logger"
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

	if host == "" || port == "" || name == "" || user == "" || pass == "" {
		log.Fatal("database configuration is incomplete; fill out all required environment variables")
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
