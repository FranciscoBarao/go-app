package controllers

import (
	"net/http"
	"user-management/middleware"
	"user-management/models"
	"user-management/services"
	"user-management/utils"

	"github.com/unrolled/render"
)

type userService interface {
	Register(user *models.User) error
	GetAll(sort string) ([]models.User, error)
	Login(username, password string) error
	Delete(name string) error
}

type UserController struct {
	service userService
}

// InitController initializes the user controller.
func InitUserController(userSvc *services.UserService) *UserController {
	return &UserController{
		service: userSvc,
	}
}

// Register User godoc
// @Summary 	Registers a User
// @Tags 		tags
// @Produce 	json
// @Param 		data body models.User true "The User name"
// @Success 	200 {object} models.User
// @Router 		/user [post]
func (controller *UserController) Register(w http.ResponseWriter, r *http.Request) {

	var user models.User
	if err := utils.DecodeJSONBody(w, r, &user); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	controller.service.Register(&user)

	render.New().JSON(w, http.StatusOK, user)
}

// Get Tags godoc
// @Summary 	Fetches all Tags
// @Tags 		tags
// @Produce 	json
// @Success 	200 {object} models.User
// @Router 		/user [get]
func (controller *UserController) GetAll(w http.ResponseWriter, r *http.Request) {

	render.New().JSON(w, http.StatusOK, "user")
}

// Get User godoc
// @Summary 	Fetches a specific User using a name
// @Tags 		tags
// @Produce 	json
// @Param 		name path string true "The User name"
// @Success 	200 {object} models.User
// @Router 		/user/{name} [get]
func (controller *UserController) Get(w http.ResponseWriter, r *http.Request) {

	render.New().JSON(w, http.StatusOK, "user")
}

// Delete User godoc
// @Summary 	Deletes a specific User
// @Tags 		tags
// @Produce 	json
// @Param 		name path string true "The User name"
// @Success 	204
// @Router 		/user/{name} [delete]
func (controller *UserController) Delete(w http.ResponseWriter, r *http.Request) {

	render.New().JSON(w, http.StatusNoContent, "user")
}
