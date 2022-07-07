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

	/*err := repo.db.Create(offer, omits...)
	if err != nil {
		return err
	}*/

	log.Println("Offer repository -> Create")

	return nil
}
