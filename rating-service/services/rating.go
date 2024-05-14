package services

import (
	"rating-service/model"
	"rating-service/repositories"
)

type ratingRepository interface {
	Create(rating *model.Rating) error
	GetAll(sort string) ([]model.Rating, error)
	Get(id string) (model.Rating, error)
	Delete(rating *model.Rating) error
}

type RatingService struct {
	repo ratingRepository
}

func InitRatingService(tagRepo *repositories.RatingRepository) *RatingService {
	return &RatingService{
		repo: tagRepo,
	}
}

func (svc *RatingService) Create(rating *model.Rating) error {

	return svc.repo.Create(rating)
}

func (svc *RatingService) GetAll(sort string) ([]model.Rating, error) {

	return svc.repo.GetAll(sort)
}

func (svc *RatingService) Get(id string) (model.Rating, error) {

	return svc.repo.Get(id)
}

func (svc *RatingService) Delete(id string) error {

	// Get rating by id
	rating, err := svc.repo.Get(id)
	if err != nil {
		return err
	}

	// Delete by id
	return svc.repo.Delete(&rating)
}
