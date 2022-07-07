package repository

import (
	"marketplace/database"
	"marketplace/repository/offerRepo"
)

// Repositories contains all the repo structs
type Repositories struct {
	OfferRepository *offerRepo.OfferRepository
}

// InitRepositories should be called in main.go
func InitRepositories(db *database.PostgresqlRepository) *Repositories {
	offerRepository := offerRepo.NewOfferRepository(db)

	return &Repositories{
		OfferRepository: offerRepository,
	}
}
