package config

import (
	"fmt"
	"os"

	"github.com/lucasmcclean/url-shortener/logger"
)

type DB struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	URL      string
}

func GetDB(log logger.Logger) *DB {
	dbCfg := &DB{}
	var missing []string

	dbCfg.URL = os.Getenv("DB_URL")
	if dbCfg.URL == "" {
		dbCfg.User, missing = getOrAppendMissing("DB_USER", missing)
		dbCfg.Password, missing = getOrAppendMissing("DB_PASSWORD", missing)
		dbCfg.Host, missing = getOrAppendMissing("DB_HOST", missing)
		dbCfg.Port, missing = getOrAppendMissing("DB_PORT", missing)
		dbCfg.Name, missing = getOrAppendMissing("DB_NAME", missing)
	}

	if len(missing) > 0 {
		log.Fatal("missing DB environment variables; must provide URL or all of these\n", "keys", missing)
	}

	return dbCfg
}

func (db *DB) GenerateConnStr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
	)
}

func (db *DB) GenerateConnStrNoSSL() string {
	if db.URL == "" {
		return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			db.User,
			db.Password,
			db.Host,
			db.Port,
			db.Name,
		)
	}
	return fmt.Sprintf("%s?sslmode=disable", db.URL)
}
