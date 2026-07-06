package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/FranciscoBarao/catalog/internal/transport"
)

// AddMechanismRouter registers mechanism routes on the given router.
func AddMechanismRouter(router chi.Router, oauthKey string, mechanismController *transport.MechanismController) {
	router.Route("/api/mechanism", func(router chi.Router) {
		router.Post("/", mechanismController.Create)
		router.Get("/", mechanismController.GetAll)
		router.Get("/{slug}", mechanismController.Get)
		router.Delete("/{slug}", mechanismController.Delete)
	})
}
