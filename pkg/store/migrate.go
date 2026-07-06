package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Migration is one embedded schema step, named NNNN_name.sql.
type Migration struct {
	Version int
	Name    string
	body    string
}

// migrate applies every embedded migration newer than the recorded schema
// version, in order, each inside its own transaction. It returns the ones it
// applied so `tama db migrate` can print them.
func (d *DB) migrate(ctx context.Context) ([]Migration, error) {
	if _, err := d.write.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at INTEGER NOT NULL
	)`); err != nil {
		return nil, err
	}

	all, err := loadMigrations()
	if err != nil {
		return nil, err
	}

	var applied []Migration
	for _, m := range all {
		var done int
		if err := d.write.QueryRowContext(ctx, `SELECT COUNT(*) FROM schema_migrations WHERE version = ?`, m.Version).Scan(&done); err != nil {
			return nil, err
		}
		if done > 0 {
			continue
		}
		err := transact(ctx, d.write, nil, func(tx *sql.Tx) error {
			if _, err := tx.ExecContext(ctx, m.body); err != nil {
				return err
			}
			_, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version, name, applied_at) VALUES (?, ?, ?)`,
				m.Version, m.Name, nowMs())
			return err
		})
		if err != nil {
			return nil, fmt.Errorf("migration %04d (%s.sql): %w", m.Version, m.Name, err)
		}
		applied = append(applied, m)
	}
	return applied, nil
}

// Applied lists the migrations this Open ran, empty when the schema was
// already current.
func (d *DB) Applied() []Migration {
	return d.applied
}

// SchemaVersion reports the highest applied migration and how many embedded
// ones are still pending, for `tama db status`.
func (d *DB) SchemaVersion(ctx context.Context) (current, pending int, err error) {
	if err := d.read.QueryRowContext(ctx, `SELECT COALESCE(MAX(version), 0) FROM schema_migrations`).Scan(&current); err != nil {
		return 0, 0, err
	}
	all, err := loadMigrations()
	if err != nil {
		return 0, 0, err
	}
	for _, m := range all {
		if m.Version > current {
			pending++
		}
	}
	return current, pending, nil
}

func loadMigrations() ([]Migration, error) {
	names, err := fs.Glob(migrations, "migrations/*.sql")
	if err != nil {
		return nil, err
	}
	sort.Strings(names)

	out := make([]Migration, 0, len(names))
	for _, name := range names {
		base := strings.TrimSuffix(path.Base(name), ".sql")
		num, rest, ok := strings.Cut(base, "_")
		if !ok {
			return nil, fmt.Errorf("migration %s: want NNNN_name.sql", name)
		}
		version, err := strconv.Atoi(num)
		if err != nil {
			return nil, fmt.Errorf("migration %s: bad version: %w", name, err)
		}
		body, err := migrations.ReadFile(name)
		if err != nil {
			return nil, err
		}
		out = append(out, Migration{Version: version, Name: rest, body: string(body)})
	}
	return out, nil
}
