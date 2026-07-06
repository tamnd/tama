package store

import "context"

// User mirrors the users table. PasswordHash is an argon2id PHC string.
type User struct {
	ID           int64
	Username     string
	PasswordHash string
	DisplayName  string
	Avatar       string
	IsAdmin      bool
	CreatedAt    int64
	SettingsJSON string
}

const userCols = "id, username, password_hash, display_name, avatar, is_admin, created_at, settings_json"

func scanUser(row interface{ Scan(...any) error }) (User, error) {
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.DisplayName, &u.Avatar, &u.IsAdmin, &u.CreatedAt, &u.SettingsJSON)
	return u, one(err)
}

const createUserSQL = `INSERT INTO users (username, password_hash, is_admin, created_at) VALUES (?, ?, ?, ?)
	RETURNING ` + userCols

// CreateUser inserts a user; the username must already be validated and
// lowercased. A duplicate surfaces as a UNIQUE constraint error.
func (d *DB) CreateUser(ctx context.Context, username, passwordHash string, isAdmin bool) (User, error) {
	return scanUser(d.write.QueryRowContext(ctx, createUserSQL, username, passwordHash, isAdmin, nowMs()))
}

const userByUsernameSQL = `SELECT ` + userCols + ` FROM users WHERE username = ?`

// UserByUsername looks a user up by name, ErrNotFound when absent.
func (d *DB) UserByUsername(ctx context.Context, username string) (User, error) {
	return scanUser(d.read.QueryRowContext(ctx, userByUsernameSQL, username))
}

const userByIDSQL = `SELECT ` + userCols + ` FROM users WHERE id = ?`

// UserByID looks a user up by id, ErrNotFound when absent.
func (d *DB) UserByID(ctx context.Context, id int64) (User, error) {
	return scanUser(d.read.QueryRowContext(ctx, userByIDSQL, id))
}

const updatePasswordSQL = `UPDATE users SET password_hash = ? WHERE id = ?`

// UpdatePassword swaps the stored hash, ErrNotFound for an unknown id.
func (d *DB) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	res, err := d.write.ExecContext(ctx, updatePasswordSQL, passwordHash, id)
	if err != nil {
		return err
	}
	return affected(res)
}

const deleteUserSQL = `DELETE FROM users WHERE id = ?`

// DeleteUser removes a user; sessions and progress cascade.
func (d *DB) DeleteUser(ctx context.Context, id int64) error {
	res, err := d.write.ExecContext(ctx, deleteUserSQL, id)
	if err != nil {
		return err
	}
	return affected(res)
}

const listUsersSQL = `SELECT ` + userCols + ` FROM users ORDER BY id`

// ListUsers returns every user, oldest first.
func (d *DB) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := d.read.QueryContext(ctx, listUsersSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
