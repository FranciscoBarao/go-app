package transport

import (
	"context"
	"net/http"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/marketplace/internal/middleware"
	"github.com/FranciscoBarao/marketplace/internal/offer"
	"github.com/FranciscoBarao/marketplace/internal/utils"
)

// OfferService defines the interface for offer business logic.
type OfferService interface {
	Create(ctx context.Context, offer *offer.Offer, username string) error
	GetAll(ctx context.Context) ([]offer.Offer, error)
	Get(ctx context.Context, uuid string) (offer.Offer, error)
	Update(ctx context.Context, req *offer.UpdateOfferRequest, uuid, username string) (offer.Offer, error)
	Delete(ctx context.Context, uuid, username string) error
}

// OfferController handles HTTP requests for offer operations.
type OfferController struct {
	service OfferService
}

// NewOfferController initializes the offer controller.
func NewOfferController(svc OfferService) *OfferController {
	return &OfferController{service: svc}
}

// Create Offer godoc
// @Summary 	Creates a Offer based on a json body
// @Tags 		offer
// @Produce 	json
// @Param 		data body offer.CreateOfferRequest true "The input Offer struct"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} offer.Offer
// @Router 		/offer [post]
func (c *OfferController) Create(w http.ResponseWriter, r *http.Request) {
	var req offer.CreateOfferRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	user, err := utils.GetUsernameFromToken(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	o := offer.NewOffer(&req)

	if err := c.service.Create(context.Background(), o, user); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, o); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// GetAll Offers godoc
// @Summary 	Fetches all Offers
// @Tags 		offer
// @Produce 	json
// @Success 	200 {object} offer.Offer
// @Router 		/offer [get]
func (c *OfferController) GetAll(w http.ResponseWriter, r *http.Request) {
	offers, err := c.service.GetAll(context.Background())
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, offers); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// Get Offer godoc
// @Summary 	Fetches a Offer
// @Tags 		offer
// @Produce 	json
// @Param 		id path string true "The Offer UUID"
// @Success 	200 {object} offer.Offer
// @Router 		/offer/{id} [get]
func (c *OfferController) Get(w http.ResponseWriter, r *http.Request) {
	uuid := utils.GetFieldFromURL(r, "id")

	o, err := c.service.Get(context.Background(), uuid)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, o); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// Update Offer godoc
// @Summary 	Updates a specific Offer via UUID
// @Tags 		offer
// @Produce 	json
// @Param 		id path string true "The Offer UUID"
// @Param 		data body offer.UpdateOfferRequest true "The Offer fields to update"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} offer.Offer
// @Router 		/offer/{id} [patch]
func (c *OfferController) Update(w http.ResponseWriter, r *http.Request) {
	var req offer.UpdateOfferRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	user, err := utils.GetUsernameFromToken(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	uuid := utils.GetFieldFromURL(r, "id")

	o, err := c.service.Update(context.Background(), &req, uuid, user)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, o); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// Delete Offer godoc
// @Summary 	Deletes a specific Offer via UUID
// @Tags 		offer
// @Produce 	json
// @Param 		id path string true "The Offer UUID"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	204
// @Router 		/offer/{id} [delete]
func (c *OfferController) Delete(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUsernameFromToken(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	uuid := utils.GetFieldFromURL(r, "id")

	if err := c.service.Delete(context.Background(), uuid, user); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusNoContent, uuid); err != nil {
		middleware.ErrorHandler(w, err)
	}
}
