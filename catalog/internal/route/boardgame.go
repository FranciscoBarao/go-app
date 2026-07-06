package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/FranciscoBarao/catalog/internal/transport"
)

// AddBoardGameRouter registers boardgame routes on the given router.
func AddBoardGameRouter(router chi.Router, oauthKey string, boardGameControler *transport.BoardgameController) {
	router.Group(func(router chi.Router) {
		router.Post("/api/boardgame", boardGameControler.Create)
		router.Post("/api/boardgame/{slug}/expansion", boardGameControler.Create)
		router.Patch("/api/boardgame/{slug}", boardGameControler.Update)
		router.Delete("/api/boardgame/{slug}", boardGameControler.Delete)
		router.Post("/api/boardgame/{slug}/rate", boardGameControler.Rate)
	})

	router.Group(func(r chi.Router) {
		router.Get("/api/boardgame", boardGameControler.GetAll)
		router.Get("/api/boardgame/by-id/{id}", boardGameControler.GetByID)
		router.Get("/api/boardgame/{slug}", boardGameControler.Get)
	})
}
