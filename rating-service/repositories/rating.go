package repositories

import (
	"errors"
	"rating-service/database"
	"rating-service/middleware"
	"rating-service/model"
)

type RatingRepository struct {
	db *database.PostgresqlRepository
}

func NewRatingRepository(instance *database.PostgresqlRepository) *RatingRepository {
	return &RatingRepository{
		db: instance,
	}
}

func (repo *RatingRepository) Create(rating *model.Rating) error {

	return repo.db.Create(rating)
}

func (repo *RatingRepository) GetAll(sort string) ([]model.Rating, error) {

	var ratings []model.Rating
	return ratings, repo.db.Read(&ratings, sort, "", "")
}

func (repo *RatingRepository) Get(id string) (model.Rating, error) {

	var rating model.Rating
	err := repo.db.Read(&rating, "", "id = ?", id)

	var mr *middleware.MalformedRequest
	if err != nil && errors.As(err, &mr) {
		return rating, middleware.NewError(mr.GetStatus(), "Rating not found with id: "+id)
	}

	return rating, err
}

func (repo *RatingRepository) Delete(rating *model.Rating) error {

	return repo.db.Delete(rating)
}
