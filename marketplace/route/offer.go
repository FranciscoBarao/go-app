package route

import (
	"marketplace/controllers"

	"github.com/go-chi/chi/v5"
)

func AddOfferRouter(router chi.Router, controller *controllers.OfferController) {
	router.Post("/api/offer", controller.Create)
	router.Get("/api/offer", controller.GetAll)
	router.Patch("/api/offer/{id}", controller.Update)
	router.Delete("/api/offer/{id}", controller.Delete)

	//router.Get("/api/offer/{id}", controller.Get)
	//router.Patch("/api/offer/{id}", controller.Update)
	//router.Delete("/api/offer/{id}", controller.Delete)
}
