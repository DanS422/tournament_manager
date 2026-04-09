package db

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

func MigrateUp() {
	m := newMigrate()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func MigrateDown() {
	m := newMigrate()

	if err := m.Steps(-1); err != nil {
		log.Fatal(err)
	}
}

func MigrateReset() {
	m := newMigrate()

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func newMigrate() *migrate.Migrate {
	wd, _ := os.Getwd()

	m, err := migrate.New(
		"file:///"+wd+"/migrations",
		"sqlite://tournaments.db",
	)

	if err != nil {
		log.Fatal(err)
	}

	return m
}
