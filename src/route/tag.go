package route

import (
	"go-app/controller/tag"

	"github.com/go-chi/chi/v5"
)

func AddTagRouter(router chi.Router, tagController *tag.Controller) {
	router.Post("/api/tag", tagController.Create)
	router.Get("/api/tag", tagController.GetAll)
	router.Get("/api/tag/{name}", tagController.Get)
	router.Delete("/api/tag/{name}", tagController.Delete)
}
