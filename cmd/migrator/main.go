package main

import (
	"errors"
	"fmt"
	"url-shortener/internal/config"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.MustLoad()

	m, err := migrate.New("file://"+cfg.MigrationPath, fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", cfg.StoragePath, cfg.MigrationTable))
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migration changes are applied")
			return
		}
		panic(err)
	}
}
