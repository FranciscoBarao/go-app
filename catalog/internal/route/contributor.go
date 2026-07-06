package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/FranciscoBarao/catalog/internal/transport"
)

// AddContributorRouter registers contributor routes.
func AddContributorRouter(router chi.Router, oauthKey string, controller *transport.ContributorController) {
	router.Route("/api/contributors", func(router chi.Router) {
		router.Post("/", controller.Create)
		router.Get("/", controller.GetAll)
		router.Get("/{slug}", controller.Get)
		router.Patch("/{slug}", controller.Update)
		router.Delete("/{slug}", controller.Delete)
	})
}
