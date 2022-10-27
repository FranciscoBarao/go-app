package repositories

import (
	"user-management/database"
)

// Repositories contains all the repo structs
type Repositories struct {
	UserRepository *UserRepository
}

// InitRepositories should be called in main.go
func InitRepositories(db *database.PostgresqlRepository) *Repositories {
	userRepository := NewUserRepository(db)

	return &Repositories{
		UserRepository: userRepository,
	}
}
