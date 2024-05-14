package controllers

import (
	"marketplace/services"
)

// Controllers contains all the controllers
type Controllers struct {
	OfferController *OfferController
}

// InitControllers returns a new Controllers
func InitControllers(services *services.Services) *Controllers {
	return &Controllers{
		OfferController: InitOfferController(services.OfferService),
	}
}
