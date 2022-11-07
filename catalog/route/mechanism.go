package route

import (
	"catalog/controllers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"
)

func AddMechanismRouter(router chi.Router, oauthKey string, mechanismController *controllers.MechanismController) {
	// Protected layer
	router.Route("/api/mechanism", func(router chi.Router) {
		// Use the Bearer Authentication middleware
		router.Use(oauth.Authorize(oauthKey, nil))

		router.Post("/", mechanismController.Create)
		router.Get("/", mechanismController.GetAll)
		router.Get("/{name}", mechanismController.Get)
		router.Delete("/{name}", mechanismController.Delete)
	})
}
