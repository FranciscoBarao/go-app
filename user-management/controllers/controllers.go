package controllers

import "user-management/services"

// Controllers contains all the controllers
type Controllers struct {
	UserController     *UserController
	VerifierController VerifierController
}

// InitControllers returns a new Controllers
func InitControllers(services *services.Services) *Controllers {
	return &Controllers{
		UserController:     InitUserController(services.UserService),
		VerifierController: InitVerifierController(services.UserService),
	}
}
