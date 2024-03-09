package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	ErrURLNotFound = errors.New("url not found") //move later?
	ErrURLExists   = errors.New("url exists")
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS url (
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	// Create the index if not exists
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
	`)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) error {
	const op = "storage.postgresql.saveURl"
	var id int
	err := s.db.QueryRow("INSERT INTO url(url, alias) VALUES ($1, $2) RETURNING id", urlToSave, alias).Scan(&id)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgresql.getURL"
	var urlToSave string
	err := s.db.QueryRow("SELECT url FROM url WHERE alias = $1", alias).Scan(&urlToSave)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s:%w", op, sql.ErrNoRows)
		}
		return "", fmt.Errorf("%s:%w", op, err)
	}
	return urlToSave, nil
}

func (s *Storage) DeleteUrl(alias string) (string, error) {
	const op = "storage.postgresql.deleteURL"
	var urlToDelete string
	err := s.db.QueryRow("DELETE FROM url WHERE alias = $1 RETURNING url", alias).Scan(&urlToDelete)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}
	return urlToDelete, nil
}
