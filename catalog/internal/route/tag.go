package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/FranciscoBarao/catalog/internal/transport"
)

// AddTagRouter registers tag routes on the given router.
func AddTagRouter(router chi.Router, oauthKey string, tagController *transport.TagController) {
	// Protected layer
	router.Route("/api/tag", func(router chi.Router) {
		// Use the Bearer Authentication middleware
		//router.Use(oauth.Authorize(oauthKey, nil))

		router.Post("/", tagController.Create)
		router.Get("/", tagController.GetAll)
		router.Get("/{name}", tagController.Get)
		router.Delete("/{name}", tagController.Delete)
	})
}
