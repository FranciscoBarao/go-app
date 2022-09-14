package offer

import (
	"log"
	"net/http"

	"marketplace/model"
	"marketplace/repository/offerRepo"
	"marketplace/utils"
)

// Declaring the repository interface in the controller package allows us to easily swap out the actual implementation, enforcing loose coupling.
type repository interface {
	Create(offer *model.Offer) error
	ReadAll() ([]model.Offer, error)
}

// Controller contains the service, which contains database-related logic, as an injectable dependency, allowing us to decouple business logic from db logic.
type Controller struct {
	repo repository
}

// InitController initializes the boargame and the associations controller.
func InitController(offerRepo *offerRepo.OfferRepository) *Controller {
	return &Controller{
		repo: offerRepo,
	}
}

// Create Offer godoc
// @Summary 	Creates a Offer based on a json body
// @Tags 		offers
// @Produce 	json
// @Param 		data body model.Offer true "The input Offer struct"
// @Success 	200 {object} model.Offer
// @Router 		/offer [post]
func (controller *Controller) Create(w http.ResponseWriter, r *http.Request) {

	log.Println(" Offer Controller -> Create ")

	var offer model.Offer
	err := utils.DecodeJSONBody(w, r, &offer)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	err = controller.repo.Create(&offer)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &offer, http.StatusOK, nil)
}

// Get Offer godoc
// @Summary 	Fetches all Offers
// @Tags 		offer
// @Produce 	json
// @Success 	200 {object} model.Offer
// @Router 		/offer [get]
func (controller *Controller) GetAll(w http.ResponseWriter, r *http.Request) {

	offers, err := controller.repo.ReadAll()
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &offers, http.StatusOK, nil)
}
