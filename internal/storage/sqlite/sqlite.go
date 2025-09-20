package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	/*stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS url(
	        					 id INTEGER PRIMARY KEY,
	        					 alias TEXT NOT NULL UNIQUE,
	        					 url TEXT NOT NULL);
	    						 CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}
		defer stmt.Close()

		_, err = stmt.Exec()
		if err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}*/

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(alias string, url string) (int64, error) {
	const op = "sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url (alias,url) VALUES (?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %v", op, err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(alias, url)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %v", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %v", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias=?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %v", op, err)
	}
	defer stmt.Close()

	var url string //-----------------------------------------Переменная для сохранения результата
	if err = stmt.QueryRow(alias).Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: execute statement: %v", op, err)
	}
	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM URL WHERE alias=?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %v", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrNoExtended(sqlite3.ErrNotFound) {
			return storage.ErrURLNotFound
		}
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}
