package route

import (
	"catalog/controllers"

	"github.com/go-chi/chi/v5"
)

func AddCategoryRouter(router chi.Router, categoryController *controllers.CategoryController) {
	router.Post("/api/category", categoryController.Create)
	router.Get("/api/category", categoryController.GetAll)
	router.Get("/api/category/{name}", categoryController.Get)
	router.Delete("/api/category/{name}", categoryController.Delete)
}
