package user

import (
	"context"

	"github.com/FranciscoBarao/user-management/internal/logging"
)

// Database defines the persistence operations needed by the user service.
type Database interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	DeleteUser(ctx context.Context, username string) error
}

// Service holds the database dependency and implements user business logic.
type Service struct {
	db Database
}

// NewService creates a new user Service.
func NewService(db Database) *Service {
	return &Service{db: db}
}

// Register hashes the password and persists a new user.
func (svc *Service) Register(ctx context.Context, user *User, password string) error {
	logging.FromCtx(ctx).Debug().Str("username", user.Username).Msg("registering user")

	if err := user.HashPassword(password); err != nil {
		return err
	}

	return svc.db.CreateUser(ctx, user)
}

// Login validates credentials by fetching the user and checking the password.
func (svc *Service) Login(ctx context.Context, username, password string) error {
	user, err := svc.db.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	return user.CheckPassword(password)
}

// GetAll retrieves all active users.
func (svc *Service) GetAll(ctx context.Context) ([]User, error) {
	return svc.db.GetAllUsers(ctx)
}

// Delete soft-deletes a user by username.
func (svc *Service) Delete(ctx context.Context, username string) error {
	logging.FromCtx(ctx).Debug().Str("username", username).Msg("deleting user")
	return svc.db.DeleteUser(ctx, username)
}
