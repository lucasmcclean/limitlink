package link

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/lucasmcclean/url-shortener/internal/config"
	"github.com/lucasmcclean/url-shortener/internal/logger"
)

type LinkRepository interface {
	Close() error
	Migrate(log logger.MigrateLogger) error

	CreateLink(original, short string, maxUses *int, expiresAt *time.Time) (uuid.UUID, error)
	GetByShort(ctx context.Context, short string) (Link, error)
	GetByAdminToken(ctx context.Context, token uuid.UUID) (Link, error)
	Visit(ctx context.Context, short string) (string, error)
	DeleteByAdminToken(ctx context.Context, token uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type Postgres struct {
	pool    *sql.DB
	connStr string
}

func NewRepository(cfg *config.DB) (LinkRepository, error) {
	connStr := cfg.GenerateConnStrNoSSL()

	pool, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %s", err)
	}

	err = pool.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging database: %s", err)
	}

	return &Postgres{
		pool:    pool,
		connStr: connStr,
	}, nil
}

func (db *Postgres) Close() error {
	return db.pool.Close()
}

func (db *Postgres) Migrate(log logger.MigrateLogger) error {
	driver, err := postgres.WithInstance(db.pool, &postgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
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

// CreateLink inserts a new shortened link with optional constraints.
func (db *Postgres) CreateLink(original, short string, maxUses *int, expiresAt *time.Time) (uuid.UUID, error) {
	adminToken := uuid.New()

	_, err := db.pool.Exec(`
		INSERT INTO link (short, original, max_uses, expires_at, admin_token)
		VALUES ($1, $2, $3, $4, $5)
	`, short, original, maxUses, expiresAt, adminToken)

	return adminToken, err
}

// GetByShort retrieves a link by its short URL.
func (db *Postgres) GetByShort(ctx context.Context, short string) (Link, error) {
	var link Link
	err := db.pool.QueryRowContext(ctx, `
		SELECT id, original, short, admin_token, max_uses, click_count, expires_at, created_at
		FROM link
		WHERE short = $1
	`, short).Scan(
		&link.ID,
		&link.Original,
		&link.Short,
		&link.AdminToken,
		&link.MaxUses,
		&link.ClickCount,
		&link.ExpiresAt,
		&link.CreatedAt,
	)
	return link, err
}

// GetByAdminToken retrieves a link by its admin token.
func (db *Postgres) GetByAdminToken(ctx context.Context, token uuid.UUID) (Link, error) {
	var link Link
	err := db.pool.QueryRowContext(ctx, `
		SELECT id, original, short, admin_token, max_uses, click_count, expires_at, created_at
		FROM link
		WHERE admin_token = $1
	`, token).Scan(
		&link.ID,
		&link.Original,
		&link.Short,
		&link.AdminToken,
		&link.MaxUses,
		&link.ClickCount,
		&link.ExpiresAt,
		&link.CreatedAt,
	)
	return link, err
}

// Visit increments usage and returns the original URL.
func (db *Postgres) Visit(ctx context.Context, short string) (string, error) {
	var url string
	err := db.pool.QueryRowContext(ctx, `
		WITH candidate AS (
			SELECT * FROM link
			WHERE short = $1
			AND (expires_at IS NULL OR expires_at > now())
		),
		updated AS (
			UPDATE link
			SET click_count = click_count + 1
			WHERE short IN (SELECT short FROM candidate)
			AND (max_uses IS NULL OR click_count < max_uses)
			RETURNING original
		)
		DELETE FROM link
		WHERE short IN (
			SELECT short FROM updated
			WHERE max_uses IS NOT NULL AND click_count >= max_uses
		)
		SELECT original FROM updated
	`, short).Scan(&url)

	return url, err
}

// DeleteByAdminToken removes a link by its admin token.
func (db *Postgres) DeleteByAdminToken(ctx context.Context, token uuid.UUID) error {
	_, err := db.pool.ExecContext(ctx, `
		DELETE FROM link
		WHERE admin_token = $1
	`, token)
	return err
}

// DeleteExpired removes all expired links.
func (db *Postgres) DeleteExpired(ctx context.Context) error {
	_, err := db.pool.ExecContext(ctx, `
		DELETE FROM link
		WHERE expires_at IS NOT NULL AND expires_at <= now()
	`)
	return err
}
