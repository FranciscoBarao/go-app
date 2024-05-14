package services

import "marketplace/repositories"

// Repositories contains all the repo structs
type Services struct {
	OfferService *OfferService
}

// InitRepositories should be called in main.go
func InitServices(repositories *repositories.Repositories) *Services {
	offerService := InitOfferService(repositories.OfferRepository)

	return &Services{
		OfferService: offerService,
	}
}
