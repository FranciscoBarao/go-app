package route

import (
	"catalog/controller/category"

	"github.com/go-chi/chi/v5"
)

func AddCategoryRouter(router chi.Router, categoryController *category.Controller) {
	router.Post("/api/category", categoryController.Create)
	router.Get("/api/category", categoryController.GetAll)
	router.Get("/api/category/{name}", categoryController.Get)
	router.Delete("/api/category/{name}", categoryController.Delete)
}
