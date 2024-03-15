package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"

	"github.com/FranciscoBarao/catalog/controllers"
)

func AddCategoryRouter(router chi.Router, oauthKey string, categoryController *controllers.CategoryController) {
	// Protected layer
	router.Route("/api/category", func(router chi.Router) {
		// Use the Bearer Authentication middleware
		router.Use(oauth.Authorize(oauthKey, nil))

		router.Post("/", categoryController.Create)
		router.Get("/", categoryController.GetAll)
		router.Get("/{name}", categoryController.Get)
		router.Delete("/{name}", categoryController.Delete)
	})
}
