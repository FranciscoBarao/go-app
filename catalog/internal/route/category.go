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
		// MethodFunc because chi has no Query helper: QUERY is not one of its
		// built-in methods and relies on chi.RegisterMethod in main.go.
		// TODO: switch to router.Query once chi supports QUERY natively.
		router.MethodFunc("QUERY", "/", categoryController.Query)
		router.Get("/{slug}", categoryController.Get)
		router.Delete("/{slug}", categoryController.Delete)
	})
}
