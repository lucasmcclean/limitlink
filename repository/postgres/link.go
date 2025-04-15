package postgres

func (db *DB) CreateLink(original, short string, maxUses int) error {
	_, err := db.pool.Exec(`
      INSERT INTO link (short, original, max_uses)
      VALUES ($1, $2, $3)
    `, short, original, maxUses)
	return err
}

func (db *DB) GetLink(short string) (string, error) {
	var url string
	err := db.pool.QueryRow(`
      SELECT original FROM link WHERE short = $1
    `, short).Scan(&url)
	return url, err
}

func (db *DB) DeleteLink(short string) error {
	_, err := db.pool.Exec(`
	    DELETE FROM link WHERE short = $1
    `, short)
	return err
}

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
          RETURNING short, original, click_count, max_uses
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

func (db *DB) DeleteExpiredLinks() error {
	_, err := db.pool.Exec(`
      DELETE FROM link WHERE expires_at IS NOT NULL AND expires_at <= now();
    `)
  return err
}
