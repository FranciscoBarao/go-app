package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"

	"github.com/FranciscoBarao/marketplace/internal/transport"
)

// AddOfferRouter registers offer routes on the given router.
func AddOfferRouter(router chi.Router, oauthKey string, controller *transport.OfferController) {
	// Protected layer
	router.Group(func(r chi.Router) {
		r.Use(oauth.Authorize(oauthKey, nil))

		r.Post("/api/offer", controller.Create)
		r.Patch("/api/offer/{id}", controller.Update)
		r.Delete("/api/offer/{id}", controller.Delete)
	})

	// Public layer
	router.Group(func(r chi.Router) {
		r.Get("/api/offer", controller.GetAll)
		r.Get("/api/offer/{id}", controller.Get)
	})
}
