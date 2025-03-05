package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/BioMihanoid/URLShortener/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(configPath string) (*Storage, error) {
	const op = "sqlite.Storage.NewStorage"
	db, err := sql.Open("sqlite3", configPath)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	return &Storage{db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const op = "sqlite.Storage.SaveUrl"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			return fmt.Errorf("%s %w", op, storage.ErrURLExists)
		}
		return fmt.Errorf("%s %w", op, err)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "sqlite.Storage.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias=?")
	if err != nil {
		return "", fmt.Errorf("failed to get url: %w", storage.ErrURLNotFound)
	}

	var resUrl string
	err = stmt.QueryRow(alias).Scan(&resUrl)
	if err != nil {
		return "", fmt.Errorf("failed to get url: %w", err)
	}

	return resUrl, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "sqlite.Storage.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias=?")
	if err != nil {
		return fmt.Errorf("%s %w", op, storage.ErrURLNotFound)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	return nil
}
