package controllers

import (
	"log"
	"net/http"

	"marketplace/middleware"
	"marketplace/model"
	"marketplace/services"
	"marketplace/utils"

	"github.com/unrolled/render"
)

// Declaring the repository interface in the controller package allows us to easily swap out the actual implementation, enforcing loose coupling.
type offerService interface {
	Create(offer *model.Offer) error
	ReadAll() ([]model.Offer, error)
}

// OfferController contains the service, which contains database-related logic, as an injectable dependency, allowing us to decouple business logic from db logic.
type OfferController struct {
	service offerService
}

// InitController initializes the boargame and the associations controller.
func InitOfferController(offerService *services.OfferService) *OfferController {
	return &OfferController{
		service: offerService,
	}
}

// Create Offer godoc
// @Summary 	Creates a Offer based on a json body
// @Tags 		offers
// @Produce 	json
// @Param 		data body model.Offer true "The input Offer struct"
// @Success 	200 {object} model.Offer
// @Router 		/offer [post]
func (controller *OfferController) Create(w http.ResponseWriter, r *http.Request) {

	log.Println(" Offer OfferController -> Create ")

	var offer model.Offer
	if err := utils.DecodeJSONBody(w, r, &offer); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.service.Create(&offer); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, offer)
}

// Get Offer godoc
// @Summary 	Fetches all Offers
// @Tags 		offer
// @Produce 	json
// @Success 	200 {object} model.Offer
// @Router 		/offer [get]
func (controller *OfferController) GetAll(w http.ResponseWriter, r *http.Request) {

	offers, err := controller.service.ReadAll()
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, offers)
}
