package route

import (
	"catalog/controllers"

	"github.com/go-chi/chi/v5"
)

func AddBoardGameRouter(router chi.Router, boardGameControler *controllers.BoardgameController) {
	router.Post("/api/boardgame", boardGameControler.Create)
	router.Get("/api/boardgame", boardGameControler.GetAll)
	router.Get("/api/boardgame/{id}", boardGameControler.Get)
	router.Patch("/api/boardgame/{id}", boardGameControler.Update)
	router.Delete("/api/boardgame/{id}", boardGameControler.Delete)

	router.Post("/api/boardgame/{id}/expansion", boardGameControler.Create)
}
