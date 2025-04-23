package postgres

import (
	"time"

	"github.com/google/uuid"
)

// Link represents a shortened URL and its constraints, defined in 0001_create_link_table.up.sql
type Link struct {
	ID         uuid.UUID  `db:"id"`
	Original   string     `db:"original"`
	Short      string     `db:"short"`
	AdminToken uuid.UUID  `db:"admin_token"`
	MaxUses    *int       `db:"max_uses"`
	ClickCount int        `db:"click_count"`
	ExpiresAt  *time.Time `db:"expires_at"`
	CreatedAt  time.Time  `db:"created_at"`
}

// CreateLink inserts a new shortened link with optional constraints.
func (db *DB) CreateLink(original, short string, maxUses *int, expiresAt *time.Time) (uuid.UUID, error) {
	adminToken := uuid.New()

	_, err := db.pool.Exec(`
    INSERT INTO link (short, original, max_uses, expires_at, admin_token)
    VALUES ($1, $2, $3, $4, $5)
    `, short, original, maxUses, expiresAt, adminToken)

	return adminToken, err
}

// GetLinkByAdminToken retrieves a Link by its admin token.
func (db *DB) GetLinkByAdminToken(adminToken uuid.UUID) (Link, error) {
	var link Link
	err := db.pool.QueryRow(`
    SELECT id, original, short, admin_token, max_uses, click_count, expires_at, created_at
    FROM link
    WHERE admin_token = $1
	`, adminToken).Scan(
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

// DeleteLinkByAdminToken removes a link by its admin token.
func (db *DB) DeleteLinkByAdminToken(adminToken uuid.UUID) error {
	_, err := db.pool.Exec(`
    DELETE FROM link
    WHERE admin_token = $1
	`, adminToken)
	return err
}

// VisitLink increments usage if valid and returns the original URL.
func (db *DB) VisitLink(short string) (string, error) {
	var url string
	err := db.pool.QueryRow(`
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
      RETURNING original, click_count, max_uses
    ),
    deleted AS (
      DELETE FROM link
      WHERE short IN (
        SELECT short FROM updated
        WHERE max_uses IS NOT NULL AND click_count >= max_uses
      )
    )
    SELECT original FROM updated
  `, short).Scan(&url)

	return url, err
}

// DeleteExpiredLinks removes all links that have expired.
func (db *DB) DeleteExpiredLinks() error {
	_, err := db.pool.Exec(`
    DELETE FROM link
    WHERE expires_at IS NOT NULL AND expires_at <= now()
	`)
	return err
}
