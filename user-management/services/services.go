package services

import "user-management/repositories"

// Repositories contains all the repo structs
type Services struct {
	UserService     *UserService
	VerifierService VerifierService
}

// InitRepositories should be called in main.go
func InitServices(repositories *repositories.Repositories) *Services {
	userService := InitUserService(repositories.UserRepository)
	verifierService := InitVerifierService(repositories.UserRepository)

	return &Services{
		UserService:     userService,
		VerifierService: verifierService,
	}
}
