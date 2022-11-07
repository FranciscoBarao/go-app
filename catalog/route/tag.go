package route

import (
	"catalog/controllers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"
)

func AddTagRouter(router chi.Router, oauthKey string, tagController *controllers.TagController) {
	// Protected layer
	router.Route("/api/tag", func(router chi.Router) {
		// Use the Bearer Authentication middleware
		router.Use(oauth.Authorize(oauthKey, nil))

		router.Post("/", tagController.Create)
		router.Get("/", tagController.GetAll)
		router.Get("/{name}", tagController.Get)
		router.Delete("/{name}", tagController.Delete)
	})
}
