package transport

import (
	"context"
	"net/http"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/user-management/internal/middleware"
	"github.com/FranciscoBarao/user-management/internal/user"
	"github.com/FranciscoBarao/user-management/internal/utils"
)

// UserService defines the interface for user business logic.
type UserService interface {
	Register(ctx context.Context, u *user.User, password string) error
	GetAll(ctx context.Context) ([]user.User, error)
	Delete(ctx context.Context, username string) error
}

// UserController handles HTTP requests for user operations.
type UserController struct {
	service UserService
}

// NewUserController initializes the user controller.
func NewUserController(svc UserService) *UserController {
	return &UserController{service: svc}
}

// Register godoc
// @Summary 	Registers a User
// @Tags 		user
// @Produce 	json
// @Param 		data body user.RegisterRequest true "The registration payload"
// @Success 	200 {object} user.User
// @Router 		/register [post]
func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var req user.RegisterRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	u := user.NewUser(&req)

	if err := c.service.Register(r.Context(), u, req.Password); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusCreated, u); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// GetAll godoc
// @Summary 	Fetches all Users
// @Tags 		user
// @Produce 	json
// @Success 	200 {object} user.User
// @Router 		/user [get]
func (c *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.GetAll(r.Context())
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, users); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// Delete godoc
// @Summary 	Deletes a User by username
// @Tags 		user
// @Produce 	json
// @Param 		username path string true "The username"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	204
// @Router 		/user/{username} [delete]
func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	req := user.DeleteRequest{Username: utils.GetFieldFromURL(r, "username")}

	if err := utils.ValidateStruct(&req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := c.service.Delete(r.Context(), req.Username); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusNoContent, req.Username); err != nil {
		middleware.ErrorHandler(w, err)
	}
}
