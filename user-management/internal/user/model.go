package user

import (
	"net/http"
	"time"

	"github.com/FranciscoBarao/user-management/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user entity.
type User struct {
	ID        uint       `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	Username  string     `json:"username" db:"username"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"-" db:"password"`
}

// HashPassword hashes the user's password using bcrypt.
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword compares the provided password against the stored hash.
func (u *User) CheckPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return middleware.NewError(http.StatusUnauthorized, "Error - Incorrect Credentials")
	}
	return nil
}
