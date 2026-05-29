package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Radixen-Dev/handoff/internal/model"

	_ "modernc.org/sqlite"
)

// DB wraps the SQLite connection.
type DB struct {
	conn *sql.DB
}

// Open opens (or creates) the database at the given path and ensures the schema exists.
func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if _, err := conn.Exec(schema); err != nil {
		conn.Close()
		return nil, fmt.Errorf("init schema: %w", err)
	}

	d := &DB{conn: conn}
	d.gc()
	return d, nil
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.conn.Close()
}

// gc removes expired packages.
func (d *DB) gc() {
	d.conn.Exec("DELETE FROM packages WHERE expires_at < ?", time.Now().UTC())
}

// Store inserts a new knowledge package.
func (d *DB) Store(pkg *model.Package) error {
	d.gc()

	tags, _ := json.Marshal(pkg.Tags)
	_, err := d.conn.Exec(
		`INSERT INTO packages (id, name, summary, content, tags, project, created_at, expires_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		pkg.ID, pkg.Name, pkg.Summary, pkg.Content, string(tags), pkg.Project,
		pkg.CreatedAt.UTC(), pkg.ExpiresAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("store package: %w", err)
	}
	return nil
}

// GetByID retrieves a package by its ID.
func (d *DB) GetByID(id string) (*model.Package, error) {
	d.gc()
	return d.scanOne("SELECT id, name, summary, content, tags, project, created_at, expires_at FROM packages WHERE id = ?", id)
}

// GetByName retrieves the most recent package with the given name.
func (d *DB) GetByName(name string) (*model.Package, error) {
	d.gc()
	return d.scanOne("SELECT id, name, summary, content, tags, project, created_at, expires_at FROM packages WHERE name = ? ORDER BY created_at DESC LIMIT 1", name)
}

// List returns all non-expired packages, optionally filtered by project.
func (d *DB) List(project string) ([]model.Package, error) {
	d.gc()

	query := "SELECT id, name, summary, content, tags, project, created_at, expires_at FROM packages"
	var args []any

	if project != "" {
		query += " WHERE project = ?"
		args = append(args, project)
	}
	query += " ORDER BY created_at DESC"

	rows, err := d.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list packages: %w", err)
	}
	defer rows.Close()

	var packages []model.Package
	for rows.Next() {
		pkg, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		packages = append(packages, *pkg)
	}
	return packages, rows.Err()
}

// GC explicitly removes all expired packages and returns the count removed.
func (d *DB) GC() (int64, error) {
	res, err := d.conn.Exec("DELETE FROM packages WHERE expires_at < ?", time.Now().UTC())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (d *DB) scanOne(query string, args ...any) (*model.Package, error) {
	row := d.conn.QueryRow(query, args...)
	var pkg model.Package
	var tags string
	var createdAt, expiresAt string

	err := row.Scan(&pkg.ID, &pkg.Name, &pkg.Summary, &pkg.Content, &tags, &pkg.Project, &createdAt, &expiresAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("package not found")
	}
	if err != nil {
		return nil, fmt.Errorf("scan package: %w", err)
	}

	json.Unmarshal([]byte(tags), &pkg.Tags)
	pkg.CreatedAt = parseTime(createdAt)
	pkg.ExpiresAt = parseTime(expiresAt)

	return &pkg, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanRow(row scannable) (*model.Package, error) {
	var pkg model.Package
	var tags string
	var createdAt, expiresAt string

	err := row.Scan(&pkg.ID, &pkg.Name, &pkg.Summary, &pkg.Content, &tags, &pkg.Project, &createdAt, &expiresAt)
	if err != nil {
		return nil, fmt.Errorf("scan package: %w", err)
	}

	json.Unmarshal([]byte(tags), &pkg.Tags)
	pkg.CreatedAt = parseTime(createdAt)
	pkg.ExpiresAt = parseTime(expiresAt)

	return &pkg, nil
}

var timeFormats = []string{
	"2006-01-02T15:04:05Z",
	"2006-01-02 15:04:05+00:00",
	"2006-01-02 15:04:05",
	time.RFC3339,
}

func parseTime(s string) time.Time {
	for _, f := range timeFormats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Time{}
}
