package config

import "os"

type DB struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// TODO: Take log as argument for handling empty values
func GetDB() *DB {
	dbCfg := &DB{}
	dbCfg.User = os.Getenv("DB_USER")
	dbCfg.Password = os.Getenv("DB_PASSWORD")
	dbCfg.Host = os.Getenv("DB_HOST")
	dbCfg.Port = os.Getenv("DB_PORT")
	dbCfg.Name = os.Getenv("DB_NAME")
	return dbCfg
}
