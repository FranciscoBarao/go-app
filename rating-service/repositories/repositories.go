package repositories

import "rating-service/database"

// Repositories contains all the repo structs
type Repositories struct {
	RatingRepository *RatingRepository
}

// InitRepositories should be called in main.go
func InitRepositories(db *database.PostgresqlRepository) *Repositories {
	ratingRepository := NewRatingRepository(db)

	return &Repositories{
		RatingRepository: ratingRepository,
	}
}
