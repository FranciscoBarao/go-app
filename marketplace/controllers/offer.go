package controllers

import (
	"net/http"

	"marketplace/middleware"
	"marketplace/model"
	"marketplace/services"
	"marketplace/utils"

	"github.com/unrolled/render"
)

// Declaring the repository interface in the controller package allows us to easily swap out the actual implementation, enforcing loose coupling.
type offerService interface {
	Create(offer *model.Offer, user string) error
	ReadAll() ([]model.Offer, error)
	Get(uuid string) (model.Offer, error)
	Update(input *model.OfferUpdate, uuid, username string) (model.Offer, error)
	Delete(uuid, username string) error
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
// @Tags 		offer
// @Produce 	json
// @Param 		data body model.Offer true "The input Offer struct"
// @Success 	200 {object} model.Offer
// @Router 		/offer [post]
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
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

	// Get username from oauth Token
	user, err := utils.GetUsernameFromToken(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.service.Create(&offer, user); err != nil {
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
// @Tags 		offer
// @Produce 	json
// @Param 		id path int true "The Offer id"
// @Param 		data body model.Offer true "The Offer struct to be updated into"
// @Success 	200 {object} model.Offer
// @Router 		/offer/{id} [patch]
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (controller *OfferController) Update(w http.ResponseWriter, r *http.Request) {

	// Deserialize input
	var input model.OfferUpdate
	if err := utils.DecodeJSONBody(w, r, &input); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate input
	if err := utils.ValidateStruct(&input); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Get username from oauth Token
	user, err := utils.GetUsernameFromToken(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	uuid := utils.GetFieldFromURL(r, "id")

	// Updates Boardgame
	offer, err := controller.service.Update(&input, uuid, user)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, offer)
}

// Delete Offer by uuid godoc
// @Summary 	Deletes a specific Offer via Uuid
// @Tags 		offer
// @Produce 	json
// @Param 		id path int true "The Offer id"
// @Success 	204
// @Router 		/offer/{id} [delete]
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (controller *OfferController) Delete(w http.ResponseWriter, r *http.Request) {

	// Get username from oauth Token
	user, err := utils.GetUsernameFromToken(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	uuid := utils.GetFieldFromURL(r, "id")

	// Delete by Id
	if err := controller.service.Delete(uuid, user); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusNoContent, uuid)
}
