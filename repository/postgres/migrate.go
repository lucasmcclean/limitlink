package postgres

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/lucasmcclean/url-shortener/logger"
)

func (db *DB) Migrate(log logger.MigrateLogger) error {
	driver, err := postgres.WithInstance(db.pool, &postgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://repository/postgres/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	migration.Log = log

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
