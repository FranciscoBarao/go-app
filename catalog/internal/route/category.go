package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/FranciscoBarao/catalog/internal/transport"
)

// AddCategoryRouter registers category routes on the given router.
func AddCategoryRouter(router chi.Router, oauthKey string, categoryController *transport.CategoryController) {
	router.Route("/api/category", func(router chi.Router) {
		router.Post("/", categoryController.Create)
		router.Get("/", categoryController.GetAll)
		router.Get("/{slug}", categoryController.Get)
		router.Delete("/{slug}", categoryController.Delete)
	})
}
