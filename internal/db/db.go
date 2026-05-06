package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// DB wraps the SQL database connection.
type DB struct {
	conn *sql.DB
}

// New opens (or creates) a SQLite database at the given path and runs migrations.
func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := conn.Exec("PRAGMA journal_mode=WAL"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("set WAL mode: %w", err)
	}

	// Enable foreign keys
	if _, err := conn.Exec("PRAGMA foreign_keys=ON"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	d := &DB{conn: conn}
	if err := d.migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return d, nil
}

// NewInMemory creates an in-memory SQLite database (for testing).
func NewInMemory() (*DB, error) {
	return New(":memory:")
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.conn.Close()
}

// Conn returns the underlying sql.DB connection.
func (d *DB) Conn() *sql.DB {
	return d.conn
}

func (d *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
id                INTEGER PRIMARY KEY AUTOINCREMENT,
text              TEXT    NOT NULL,
status            TEXT    NOT NULL DEFAULT 'todo',
date              TEXT    NOT NULL,
created_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
updated_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
status_changed_at DATETIME
);

	CREATE INDEX IF NOT EXISTS idx_tasks_date   ON tasks(date);
	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	`
	_, err := d.conn.Exec(schema)
	return err
}
