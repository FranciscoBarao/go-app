package controllers

import (
	"marketplace/repositories"
)

// Controllers contains all the controllers
type Controllers struct {
	OfferController *OfferController
}

// InitControllers returns a new Controllers
func InitControllers(repositories *repositories.Repositories) *Controllers {
	return &Controllers{
		OfferController: InitOfferController(repositories.OfferRepository),
	}
}
