package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"rest_API/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	absPath, _ := filepath.Abs(storagePath)
	fmt.Println("📂 Absolute DB path:", absPath)

	if _, err := os.Stat(absPath); err != nil {
		fmt.Println("⚠️  Database file not found, it will be created.")
	}

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
CREATE TABLE IF NOT EXISTS url(
    id INTEGER PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url,alias) VALUES(?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statment: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: execute statment: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url where alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statment: %w", op, err)
	}
	defer stmt.Close()

	var resUrl string

	err = stmt.QueryRow(alias).Scan(&resUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}

	if err != nil {
		return "", fmt.Errorf("%s: execute statment: %w", op, err)
	}

	return resUrl, nil
}

func (s *Storage) GetAllURLs() (map[string]string, error) {
	const op = "storage.sqlite.GetAllURLs"

	stmt, err := s.db.Prepare("SELECT alias, url FROM url")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statment: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: execute Query: %w", op, err)
	}
	defer rows.Close()

	urls := make(map[string]string)

	for rows.Next() {
		var alias, url string
		if err := rows.Scan(&alias, &url); err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}

		urls[alias] = url
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return urls, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url where alias = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statment: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: execute statment: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: effected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows deleted - alias '%s' not found", alias)
	}

	return nil
}
