package services

import (
	"marketplace/model"
	"marketplace/repositories"
)

type offerRepository interface {
	Create(offer *model.Offer) error
	ReadAll() ([]model.Offer, error)
	Update(offer model.Offer) error
	Get(id, username string) (model.Offer, error)
	Delete(id string) error
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

func (svc *OfferService) Create(offer *model.Offer, user string) error {

	offer.SetUsername(user)

	return svc.repo.Create(offer)
}

func (svc *OfferService) ReadAll() ([]model.Offer, error) {

	offers, err := svc.repo.ReadAll()
	if err != nil {
		return offers, err
	}

	return offers, nil
}

func (svc *OfferService) Get(uuid string) (model.Offer, error) {

	// Get Offer by id
	offer, err := svc.repo.Get(uuid, "")
	if err != nil {
		return offer, err
	}

	return offer, nil
}

func (svc *OfferService) Update(input *model.OfferUpdate, uuid, username string) (model.Offer, error) {

	// Get Offer by id & username
	offer, err := svc.repo.Get(uuid, username)
	if err != nil {
		return model.Offer{}, err
	}

	offer.UpdateOffer(input)

	return offer, svc.repo.Update(offer)
}

func (svc *OfferService) Delete(id, username string) error {

	// Get Offer by id
	offer, err := svc.repo.Get(id, username)
	if err != nil {
		return err
	}

	return svc.repo.Delete(offer.GetId())
}
