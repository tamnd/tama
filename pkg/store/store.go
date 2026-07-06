// Package store owns the SQLite database: opening it with sane pragmas,
// running the embedded migrations, and giving the rest of the app typed
// queries instead of raw SQL sprinkled everywhere. It is the only package
// that imports database/sql.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	_ "modernc.org/sqlite"
)

// Now is the clock for every timestamp the store writes. Tests swap it to
// freeze time for expiry and due-queue assertions.
var Now = func() time.Time { return time.Now().UTC() }

// nowMs is Now as UTC unix milliseconds, the unit of every timestamp column.
func nowMs() int64 {
	return Now().UnixMilli()
}

// ErrNotFound is returned by single-row lookups that match nothing.
var ErrNotFound = errors.New("store: not found")

// DB wraps two pools over one SQLite file: a read pool sized to GOMAXPROCS
// and a write pool with exactly one connection, so writers queue in Go
// instead of fighting over SQLITE_BUSY.
type DB struct {
	read  *sql.DB
	write *sql.DB

	applied []Migration
}

const pragmas = "?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=busy_timeout(5000)"

// Open opens (or creates) the database at path, applies the pragmas on every
// connection, and brings the schema up to date. The parent directory is
// created 0700 if missing and the db file is clamped to 0600.
func Open(ctx context.Context, path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}

	write, err := sql.Open("sqlite", path+pragmas+"&_txlock=immediate")
	if err != nil {
		return nil, err
	}
	write.SetMaxOpenConns(1)

	read, err := sql.Open("sqlite", path+pragmas)
	if err != nil {
		write.Close()
		return nil, err
	}
	read.SetMaxOpenConns(runtime.GOMAXPROCS(0))

	db := &DB{read: read, write: write}
	if db.applied, err = db.migrate(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	if err := os.Chmod(path, 0o600); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// Close checkpoints the WAL back into the main file and closes both pools.
func (d *DB) Close() error {
	_, err := d.write.Exec("PRAGMA wal_checkpoint(TRUNCATE)")
	return errors.Join(err, d.read.Close(), d.write.Close())
}

// Read runs fn inside a read-only transaction on the read pool.
func (d *DB) Read(ctx context.Context, fn func(*sql.Tx) error) error {
	return transact(ctx, d.read, &sql.TxOptions{ReadOnly: true}, fn)
}

// Write runs fn inside an immediate write transaction on the single-writer
// pool, committing on nil and rolling back on error.
func (d *DB) Write(ctx context.Context, fn func(*sql.Tx) error) error {
	return transact(ctx, d.write, nil, fn)
}

func transact(ctx context.Context, db *sql.DB, opts *sql.TxOptions, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// one scans a single row and maps sql.ErrNoRows to ErrNotFound.
func one(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

// affected maps a zero-row UPDATE or DELETE to ErrNotFound.
func affected(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
