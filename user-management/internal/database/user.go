package database

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"

	dbsql "github.com/FranciscoBarao/user-management/internal/database/sql"
	"github.com/FranciscoBarao/user-management/internal/middleware"
	"github.com/FranciscoBarao/user-management/internal/user"
)

// CreateUser inserts a new user and sets the generated ID.
func (p *Postgres) CreateUser(ctx context.Context, u *user.User) error {
	err := p.pool.QueryRow(ctx, dbsql.InsertUser, u.Username, u.Email, u.Password).Scan(&u.ID)
	return mapPgError(err)
}

// GetUserByUsername retrieves a single active user by username.
func (p *Postgres) GetUserByUsername(ctx context.Context, username string) (user.User, error) {
	rows, err := p.pool.Query(ctx, dbsql.SelectUserByUsername, username)
	if err != nil {
		return user.User{}, mapPgError(err)
	}
	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[user.User])
	return u, mapPgError(err)
}

// GetAllUsers retrieves all active users (without password).
func (p *Postgres) GetAllUsers(ctx context.Context) ([]user.User, error) {
	rows, err := p.pool.Query(ctx, dbsql.SelectAllUsers)
	if err != nil {
		return nil, mapPgError(err)
	}
	users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (user.User, error) {
		var u user.User
		err := row.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt, &u.Username, &u.Email)
		return u, err
	})
	return users, mapPgError(err)
}

// DeleteUser soft-deletes a user by username.
func (p *Postgres) DeleteUser(ctx context.Context, username string) error {
	cmd, err := p.pool.Exec(ctx, dbsql.SoftDeleteUser, username)
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}
	return nil
}
