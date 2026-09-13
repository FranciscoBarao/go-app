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
		// MethodFunc because chi has no Query helper: QUERY is not one of its
		// built-in methods and relies on chi.RegisterMethod in main.go.
		// TODO: switch to router.Query once chi supports QUERY natively.
		router.MethodFunc("QUERY", "/", mechanismController.Query)
		router.Get("/{slug}", mechanismController.Get)
		router.Delete("/{slug}", mechanismController.Delete)
	})
}
