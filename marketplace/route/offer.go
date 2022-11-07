package route

import (
	"marketplace/controllers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"
)

func AddOfferRouter(router chi.Router, oauthKey string, controller *controllers.OfferController) {
	// Protected layer
	router.Group(
		func(r chi.Router) {
			// Use the Bearer Authentication middleware
			r.Use(oauth.Authorize(oauthKey, nil))

			r.Post("/api/offer", controller.Create)
			r.Patch("/api/offer/{id}", controller.Update)
			r.Delete("/api/offer/{id}", controller.Delete)
		},
	)

	// Public layer
	router.Group(
		func(r chi.Router) {
			r.Get("/api/offer", controller.GetAll)
			r.Get("/api/offer/{id}", controller.Get)
		},
	)
}
