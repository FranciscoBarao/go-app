package services

import "user-management/repositories"

// Repositories contains all the repo structs
type Services struct {
	UserService *UserService
}

// InitRepositories should be called in main.go
func InitServices(repositories *repositories.Repositories) *Services {
	userService := InitUserService(repositories.UserRepository)

	return &Services{
		UserService: userService,
	}
}
