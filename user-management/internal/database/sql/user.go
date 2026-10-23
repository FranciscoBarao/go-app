package sql

// InsertUser inserts a new user and returns the generated ID.
const InsertUser = `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`

// SelectUserByUsername retrieves a single active user by username.
const SelectUserByUsername = `SELECT id, created_at, updated_at, deleted_at, username, email, password FROM users WHERE username = $1 AND deleted_at IS NULL`

// SelectAllUsers retrieves all active users (excludes password).
const SelectAllUsers = `SELECT id, created_at, updated_at, deleted_at, username, email FROM users WHERE deleted_at IS NULL`

// SoftDeleteUser sets deleted_at for a user by username.
const SoftDeleteUser = `UPDATE users SET deleted_at = NOW() WHERE username = $1 AND deleted_at IS NULL`
