package route

import (
	"marketplace/controller/offer"

	"github.com/go-chi/chi/v5"
)

func AddOfferRouter(router chi.Router, offerControler *offer.Controller) {
	router.Post("/api/offer", offerControler.Create)
	//router.Get("/api/offer", offerControler.GetAll)
	//router.Get("/api/offer/{id}", offerControler.Get)
	//router.Patch("/api/offer/{id}", offerControler.Update)
	//router.Delete("/api/offer/{id}", offerControler.Delete)
}
