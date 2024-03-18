package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/FranciscoBarao/catalog/controllers"
)

func AddBoardGameRouter(router chi.Router, oauthKey string, boardGameControler *controllers.BoardgameController) {
	// Protected layer
	router.Group(func(router chi.Router) {
		// Use the Bearer Authentication middleware
		//router.Use(oauth.Authorize(oauthKey, nil))

		router.Post("/api/boardgame", boardGameControler.Create)
		router.Patch("/api/boardgame/{id}", boardGameControler.Update)
		router.Delete("/api/boardgame/{id}", boardGameControler.Delete)
		router.Post("/api/boardgame/{id}/rate", boardGameControler.Rate)

		router.Post("/api/boardgame/{id}/expansion", boardGameControler.Create)

	})

	// Public layer
	router.Group(func(r chi.Router) {
		router.Get("/api/boardgame", boardGameControler.GetAll)
		router.Get("/api/boardgame/{id}", boardGameControler.Get)
	})
}
