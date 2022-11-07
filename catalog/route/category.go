package route

import (
	"catalog/controllers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"
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
