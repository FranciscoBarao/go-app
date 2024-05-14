package repositories

import (
	"marketplace/database"
)

// Repositories contains all the repo structs
type Repositories struct {
	OfferRepository *OfferRepository
}

// InitRepositories should be called in main.go
func InitRepositories(db *database.PostgresqlRepository) *Repositories {
	offerRepository := NewOfferRepository(db)

	return &Repositories{
		OfferRepository: offerRepository,
	}
}
