package controller

import (
	"marketplace/controller/offer"
	"marketplace/repository"
)

// Controllers contains all the controllers
type Controllers struct {
	OfferController *offer.Controller
}

// InitControllers returns a new Controllers
func InitControllers(repositories *repository.Repositories) *Controllers {
	return &Controllers{
		OfferController: offer.InitController(repositories.OfferRepository),
	}
}
