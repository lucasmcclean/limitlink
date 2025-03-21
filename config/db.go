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

	dbCfg.URL = os.Getenv("DB_URL")

	if dbCfg.URL == "" {
		log.Fatal("missing environment variable DB_URL")
	}

	return dbCfg
}

func (db *DB) GenerateConnStr() string {
	return db.URL
}

func (db *DB) GenerateConnStrNoSSL() string {
	return fmt.Sprintf("%s?sslmode=disable", db.URL)
}
