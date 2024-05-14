package services

import "rating-service/repositories"

// Repositories contains all the repo structs
type Services struct {
	RatingService *RatingService
}

// InitRepositories should be called in main.go
func InitServices(repositories *repositories.Repositories) *Services {
	ratingService := InitRatingService(repositories.RatingRepository)

	return &Services{
		RatingService: ratingService,
	}
}
