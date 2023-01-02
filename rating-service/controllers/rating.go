package controllers

import (
	"log"
	"net/http"

	"rating-service/middleware"
	"rating-service/model"
	"rating-service/services"
	"rating-service/utils"

	"github.com/unrolled/render"
)

type ratingService interface {
	Create(rating *model.Rating) error
	GetAll(sort string) ([]model.Rating, error)
	Get(id string) (model.Rating, error)
	Delete(id string) error
}

type RatingController struct {
	service ratingService
}

// InitController initializes the rating controller
func InitRatingController(ratingSvc *services.RatingService) *RatingController {
	return &RatingController{
		service: ratingSvc,
	}
}

// Create Rating godoc
// @Summary 	Creates a Rating using model
// @Tags 		ratings
// @Produce 	json
// @Param 		data body model.Rating true "The Rating model"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} model.Rating
// @Router 		/rating [post]
func (controller *RatingController) Create(w http.ResponseWriter, r *http.Request) {

	// Deserialize Rating input
	var rating model.Rating
	if err := utils.DecodeJSONBody(w, r, &rating); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Get username from oauth Token
	user, err := utils.GetUsernameFromToken(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	log.Println(user)

	// Validate Rating input
	if err := utils.ValidateStruct(&rating); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.service.Create(&rating); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, rating)
}

// Get Ratings godoc
// @Summary 	Fetches all Ratings
// @Tags 		ratings
// @Produce 	json
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} model.Rating
// @Router 		/rating [get]
func (controller *RatingController) GetAll(w http.ResponseWriter, r *http.Request) {

	ratings, err := controller.service.GetAll("")
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	render.New().JSON(w, http.StatusOK, ratings)
}

// Get Rating godoc
// @Summary 	Fetches a specific Rating using an id
// @Tags 		ratings
// @Produce 	json
// @Param 		id path string true "The Rating id"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} model.Rating
// @Router 		/rating/{id} [get]
func (controller *RatingController) Get(w http.ResponseWriter, r *http.Request) {

	id := utils.GetFieldFromURL(r, "id")

	rating, err := controller.service.Get(id)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	render.New().JSON(w, http.StatusOK, rating)
}

// Delete Rating godoc
// @Summary 	Deletes a specific Rating
// @Tags 		ratings
// @Produce 	json
// @Param 		id path string true "The Rating id"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	204
// @Router 		/rating/{id} [delete]
func (controller *RatingController) Delete(w http.ResponseWriter, r *http.Request) {

	id := utils.GetFieldFromURL(r, "id")

	// Delete by id
	if err := controller.service.Delete(id); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusNoContent, id)
}
