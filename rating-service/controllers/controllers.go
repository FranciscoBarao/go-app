package controllers

import (
	"rating-service/services"
)

// Controllers contains all the controllers
type Controllers struct {
	RatingController *RatingController
}

// InitControllers returns a new Controllers
func InitControllers(services *services.Services) *Controllers {
	return &Controllers{
		RatingController: InitRatingController(services.RatingService),
	}
}
