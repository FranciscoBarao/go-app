package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"

	"github.com/FranciscoBarao/user-management/internal/transport"
)

// AddUserRouter registers user routes on the given router.
func AddUserRouter(router chi.Router, oauthKey string, controller *transport.UserController) {
	// Public
	router.Post("/api/register", controller.Register)

	// Protected
	router.Group(func(r chi.Router) {
		r.Use(oauth.Authorize(oauthKey, nil))

		r.Get("/api/user", controller.GetAll)
		r.Delete("/api/user/{username}", controller.Delete)
	})
}
