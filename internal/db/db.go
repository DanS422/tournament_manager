package db

import (
	"database/sql"
	"embed"

	_ "modernc.org/sqlite"
)

// embed all migration files
//
//go:embed migrations/*.sql
var migrationFS embed.FS

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
