package offerRepo

import (
	"log"
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
	err := repo.db.Create(query, offer.GetName())
	if err != nil {
		return err
	}

	return nil
}

func (repo *OfferRepository) ReadAll() ([]model.Offer, error) {

	var offers []model.Offer

	query := "SELECT name FROM offer"
	err := repo.db.ReadAll(query, &offers)
	if err != nil {
		return nil, err
	}

	// Printing offers
	for k, v := range offers {
		log.Println(k, v)
	}

	return offers, nil
}
