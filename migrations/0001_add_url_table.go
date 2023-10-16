package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up, Down)
}

func Up(context context.Context, tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS url(
	    id SERIAL PRIMARY KEY,
	    alias TEXT NOT NULL UNIQUE,
	    url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`

	_, err := tx.ExecContext(context, query)
	if err != nil {
		return err
	}

	return nil
}

func Down(context context.Context, tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS url;
	DROP INDEX IF EXISTS idx_alias`

	_, err := tx.ExecContext(context, query)
	if err != nil {
		return err
	}

	return nil
}
