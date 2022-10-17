package repositories

import (
	"marketplace/database"
	"marketplace/model"
)

type OfferRepository struct {
	db *database.PostgresqlRepository
}

func NewOfferRepository(instance *database.PostgresqlRepository) *OfferRepository {
	return &OfferRepository{
		db: instance,
	}
}

func (repo *OfferRepository) Create(offer *model.Offer) error {

	query := `INSERT INTO offer (name) VALUES ($1)`
	return repo.db.Create(query, offer.GetName())
}

func (repo *OfferRepository) ReadAll() ([]model.Offer, error) {

	var offers []model.Offer
	query := "SELECT name FROM offer"
	return offers, repo.db.ReadAll(query, &offers)
}
