package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/pressly/goose/v3"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"rest_url_shorter/internal/config"
	"rest_url_shorter/internal/storage"
)

type Storage struct {
	db *sqlx.DB
}

func New(cfgDB config.ConnectionDB) (*Storage, error) {
	const op = "storage.postgres.New"

	// open database connection
	db, err := sqlx.Open("pgx", genURLFromConfig(cfgDB))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// healthcheck
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: Ping is error: %w", op, err)
	}

	// make migrations
	err = goose.Up(db.DB, "migrations")
	if err != nil {
		return nil, fmt.Errorf("%s: Migrations error: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func genURLFromConfig(cfg config.ConnectionDB) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Address,
		cfg.DBName,
		cfg.SSLMode,
	)
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const op = "storage.postgres.SaveURL"

	stmt, err := s.db.Prepare(`INSERT INTO url(url, alias) VALUES($1, $2)`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Printf("query params: %s , %s", urlToSave, alias)

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"
	var urlFromDB string

	err := s.db.QueryRow("SELECT url FROM url WHERE alias = $1", alias).Scan(&urlFromDB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return urlFromDB, nil
}

func (s *Storage) DeleteURL(alias string) (int64, error) {
	const op = "storage.postgres.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = $1")
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	execRes, err := stmt.Exec(alias)
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}
	affectedRows, err := execRes.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	return affectedRows, nil
}
