package route

import (
	"catalog/controllers"

	"github.com/go-chi/chi/v5"
)

func AddMechanismRouter(router chi.Router, mechanismController *controllers.MechanismController) {
	router.Post("/api/mechanism", mechanismController.Create)
	router.Get("/api/mechanism", mechanismController.GetAll)
	router.Get("/api/mechanism/{name}", mechanismController.Get)
	router.Delete("/api/mechanism/{name}", mechanismController.Delete)
}
