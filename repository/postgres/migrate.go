package postgres

import (
	"github.com/golang-migrate/migrate"
	"github.com/lucasmcclean/url-shortener/logger"
)

// Performs and logs database migrations for a Postgres database. Migrations
// are listed in repository/postgres/migrations.
func (db *DB) Migrate(log logger.MigrateLogger) error {
	migration, err := migrate.New("file://repository/postgres/migrations", db.connStr)
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
