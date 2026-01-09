package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	connStr := getTestDBURL()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping test database: %v", err)
	}

	// Reset database
	if err := DropAllTables(db); err != nil {
		log.Fatalf("failed to drop tables: %v", err)
	}
	if err := DropAllTypes(db); err != nil {
		log.Fatalf("failed to drop types: %v", err)
	}

	if err := Migrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	return db
}

func StartPostgres(ctx context.Context) (*sql.DB, func()) {
	connStr := getTestDBURL()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	cleanup := func() {
		if err := DropAllTables(db); err != nil {
			log.Printf("failed to drop tables during cleanup: %v", err)
		}
		if err := DropAllTypes(db); err != nil {
			log.Printf("failed to drop types during cleanup: %v", err)
		}
		db.Close()
	}

	return db, cleanup
}

func getTestDBURL() string {
	url := os.Getenv("TEST_DB_URL")
	if url == "" {
		log.Fatal("TEST_DB_URL is required for database tests")
	}
	return url
}

func DropAllTables(db *sql.DB) error {
	rows, err := db.Query(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_type = 'BASE TABLE'
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return err
		}
		tables = append(tables, table)
	}

	if len(tables) == 0 {
		return nil
	}

	query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", strings.Join(tables, ", "))
	_, err = db.Exec(query)
	return err
}

func DropAllTypes(db *sql.DB) error {
	rows, err := db.Query(`
		SELECT n.nspname as schema, t.typname as type
		FROM pg_type t
		LEFT JOIN pg_catalog.pg_namespace n ON n.oid = t.typnamespace
		WHERE (t.typrelid = 0 OR (SELECT c.relkind = 'c' FROM pg_catalog.pg_class c WHERE c.oid = t.typrelid))
		AND NOT EXISTS(SELECT 1 FROM pg_catalog.pg_type el WHERE el.oid = t.typelem AND el.typarray = t.oid)
		AND n.nspname = 'public'
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var types []string
	for rows.Next() {
		var schema, typ string
		if err := rows.Scan(&schema, &typ); err != nil {
			return err
		}
		types = append(types, typ)
	}

	if len(types) == 0 {
		return nil
	}

	query := fmt.Sprintf("DROP TYPE IF EXISTS %s CASCADE", strings.Join(types, ", "))
	_, err = db.Exec(query)
	return err
}

func Migrate(db *sql.DB) error {
	migrationsDir, err := migrationsPath()
	if err != nil {
		return err
	}
	return goose.Up(db, migrationsDir)
}

func TruncateTables(t *testing.T, db *sql.DB, tables ...string) {
	t.Helper()
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
		if err != nil {
			t.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}
}

func migrationsPath() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to determine migrations path")
	}
	// internal/testutil -> internal/persistance/migrations
	dir := filepath.Join(filepath.Dir(file), "..", "persistance", "migrations")
	return filepath.Abs(dir)
}
