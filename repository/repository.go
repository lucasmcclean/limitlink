package repository

import "github.com/lucasmcclean/url-shortener/logger"

type Repository interface {
	Migrate(log logger.MigrateLogger) error
	Close() error
}
