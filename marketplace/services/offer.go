package services

import (
	"marketplace/model"
	"marketplace/repositories"
)

type offerRepository interface {
	Create(offer *model.Offer) error
	ReadAll() ([]model.Offer, error)
}

// Controller contains the service, which contains database-related logic, as an injectable dependency, allowing us to decouple business logic from db logic.
type OfferService struct {
	repo offerRepository
}

// InitController initializes the boargame and the associations controller.
func InitOfferService(offerRepository *repositories.OfferRepository) *OfferService {
	return &OfferService{
		repo: offerRepository,
	}
}

func (svc *OfferService) Create(offer *model.Offer) error {

	return svc.repo.Create(offer)
}

func (svc *OfferService) ReadAll() ([]model.Offer, error) {

	offers, err := svc.repo.ReadAll()
	if err != nil {
		return offers, err
	}

	return offers, nil
}
