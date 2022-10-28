package controllers

import (
	"log"
	"net/http"

	"marketplace/middleware"
	"marketplace/model"
	"marketplace/services"
	"marketplace/utils"

	"github.com/go-chi/oauth"

	"github.com/unrolled/render"
)

// Declaring the repository interface in the controller package allows us to easily swap out the actual implementation, enforcing loose coupling.
type offerService interface {
	Create(offer *model.Offer) error
	ReadAll() ([]model.Offer, error)
	Get(uuid string) (model.Offer, error)
	Update(input *model.Offer, uuid string) error
	Delete(uuid string) error
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

	// Deserialize input
	var offer model.Offer
	if err := utils.DecodeJSONBody(w, r, &offer); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate input
	if err := utils.ValidateStruct(&offer); err != nil {
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

	log.Println(r.Context().Value(oauth.CredentialContext))
	log.Println(r.Context().Value(oauth.ClaimsContext))

	offers, err := controller.service.ReadAll()
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, offers)
}

// Get Offer godoc
// @Summary 	Fetches a Offer
// @Tags 		offer
// @Produce 	json
// @Success 	200 {object} model.Offer
// @Router 		/offer/{id} [get]
func (controller *OfferController) Get(w http.ResponseWriter, r *http.Request) {

	uuid := utils.GetFieldFromURL(r, "id")

	offer, err := controller.service.Get(uuid)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, offer)
}

// Update Offer by uuid godoc
// @Summary 	Updates a specific Offer via Uuid
// @Tags 		offers
// @Produce 	json
// @Param 		id path int true "The Offer id"
// @Param 		data body model.Offer true "The Offer struct to be updated into"
// @Success 	200 {object} model.Offer
// @Router 		/offer/{id} [patch]
func (controller *OfferController) Update(w http.ResponseWriter, r *http.Request) {

	// Deserialize input
	var input model.Offer
	if err := utils.DecodeJSONBody(w, r, &input); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate input
	if err := utils.ValidateStruct(&input); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	uuid := utils.GetFieldFromURL(r, "id")

	// Updates Boardgame
	if err := controller.service.Update(&input, uuid); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, input)
}

// Delete Offer by uuid godoc
// @Summary 	Deletes a specific Offer via Uuid
// @Tags 		offers
// @Produce 	json
// @Param 		id path int true "The Offer id"
// @Success 	204
// @Router 		/offer/{id} [delete]
func (controller *OfferController) Delete(w http.ResponseWriter, r *http.Request) {

	uuid := utils.GetFieldFromURL(r, "id")

	// Delete by Id
	if err := controller.service.Delete(uuid); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusNoContent, uuid)
}
