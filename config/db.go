package config

import (
	"fmt"
	"os"
)

type DB struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	URL      string
}

// TODO: Take log as argument for handling empty values
func GetDB() *DB {
	dbCfg := &DB{}

	dbCfg.URL = os.Getenv("DB_URL")
	if dbCfg.URL == "" {
		dbCfg.User = os.Getenv("DB_USER")
		dbCfg.Password = os.Getenv("DB_PASSWORD")
		dbCfg.Host = os.Getenv("DB_HOST")
		dbCfg.Port = os.Getenv("DB_PORT")
		dbCfg.Name = os.Getenv("DB_NAME")
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
