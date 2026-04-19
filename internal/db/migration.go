package db

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

func MigrateUp(db *sql.DB) error {
	m, err := newMigrate(db)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func MigrateDown(db *sql.DB) error {
	m, err := newMigrate(db)
	if err != nil {
		return err
	}

	if err := m.Steps(-1); err != nil {
		return err
	}

	return nil
}

func MigrateReset(db *sql.DB) error {
	m, err := newMigrate(db)
	if err != nil {
		return err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func newMigrate(db *sql.DB) (*migrate.Migrate, error) {
	sourceDriver, err := iofs.New(migrationFS, "migrations")

	if err != nil {
		return nil, err
	}

	dbDriver, err := sqlite.WithInstance(db, &sqlite.Config{})

	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		"sqlite",
		dbDriver,
	)

	if err != nil {
		return nil, err
	}

	return m, nil
}
