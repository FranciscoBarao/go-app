package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/FranciscoBarao/catalog/transport"
)

func AddMechanismRouter(router chi.Router, oauthKey string, mechanismController *transport.MechanismController) {
	// Protected layer
	router.Route("/api/mechanism", func(router chi.Router) {
		// Use the Bearer Authentication middleware
		//router.Use(oauth.Authorize(oauthKey, nil))

		router.Post("/", mechanismController.Create)
		router.Get("/", mechanismController.GetAll)
		router.Get("/{name}", mechanismController.Get)
		router.Delete("/{name}", mechanismController.Delete)
	})
}
