package route

import (
	"go-app/controller/mechanism"

	"github.com/go-chi/chi/v5"
)

func AddMechanismRouter(router chi.Router, mechanismController *mechanism.Controller) {
	router.Post("/api/mechanism", mechanismController.Create)
	router.Get("/api/mechanism", mechanismController.GetAll)
	router.Get("/api/mechanism/{name}", mechanismController.Get)
	router.Delete("/api/mechanism/{name}", mechanismController.Delete)
}
