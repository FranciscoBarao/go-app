package route

import (
	"rating-service/controllers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"
)

func AddRatingRouter(router chi.Router, oauthKey string, ratingController *controllers.RatingController) {
	// Protected layer
	router.Group(func(router chi.Router) {
		// Use the Bearer Authentication middleware
		router.Use(oauth.Authorize(oauthKey, nil))

		router.Get("/api/rating", ratingController.GetAll)
		router.Get("/api/rating/{id}", ratingController.Get)
		router.Post("/api/rating", ratingController.Create)
		//router.Patch("/api/rating/{id}", ratingController.Update)
		router.Delete("/api/rating/{id}", ratingController.Delete)
	})
}
