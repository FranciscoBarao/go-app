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

	query := `INSERT INTO offer (type, name, price) VALUES ($1, $2, $3) RETURNING uuid`
	uuid, err := repo.db.Create(query, offer.GetType(), offer.GetName(), offer.GetPrice())
	if err != nil {
		return err
	}

	offer.SetId(uuid)

	return nil
}

func (repo *OfferRepository) ReadAll() ([]model.Offer, error) {

	var offers []model.Offer
	query := "SELECT * FROM offer"
	return offers, repo.db.GetAll(query, &offers)
}

func (repo *OfferRepository) Get(uuid string) (model.Offer, error) {

	var offer model.Offer
	query := `SELECT * FROM offer WHERE uuid=$1`
	return offer, repo.db.Get(query, &offer, uuid)
}

func (repo *OfferRepository) Update(offer model.Offer) error {

	query := `UPDATE offer SET name=$1, price=$2 WHERE uuid=$3`
	return repo.db.ExecuteQuery(query, offer.GetName(), offer.GetPrice(), offer.GetId())
}

func (repo *OfferRepository) Delete(uuid string) error {

	query := `DELETE FROM offer WHERE uuid=$1`
	return repo.db.ExecuteQuery(query, uuid)
}
