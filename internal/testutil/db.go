package testutil

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"tournament_manager/internal/db"
)

func SetupTestRepo(t *testing.T) *sql.DB {
	t.Helper()

	// unique DB per test to avoid cross-test mess
	// override db configuration
	dsn := fmt.Sprintf("file:test-%d.db?mode=memory&cache=shared", time.Now().UnixNano())
	dbConn, err := db.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		dbConn.Close()
	})

	if err := db.MigrateUp(dbConn); err != nil {
		t.Fatal(err)
	}

	return dbConn
}
