package boardgame

import (
	"errors"
	"net/http"

	"github.com/unrolled/render"

	"go-app/model"
	"go-app/repository/boardgameRepo"
	"go-app/utils"
)

// declaring the repository interface in the controller package allows us to easily swap out the actual implementation, enforcing loose coupling.
type repository interface {
	Create(boardgame model.BoardGame) error
	GetAll() ([]model.BoardGame, error)
	GetByName(name string) (model.BoardGame, error)
	GetById(id string) (model.BoardGame, error)
	Update(boardgame model.BoardGame) error
	DeleteById(boardgame model.BoardGame) error
}

// Controller contains the service, which contains database-related logic, as an injectable dependency, allowing us to decouple business logic from db logic.
type Controller struct {
	repo repository
}

// InitController initializes the boargame controller.
func InitController(boardGameRepo *boardgameRepo.BoardGameRepository) *Controller {
	return &Controller{
		repo: boardGameRepo,
	}
}

// Method that Creates a boardgame based on json input
func (controller *Controller) Create(w http.ResponseWriter, r *http.Request) {

	var boardGame model.BoardGame
	err := utils.DecodeJSONBody(w, r, &boardGame)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = controller.repo.Create(boardGame)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	render.New().JSON(w, http.StatusOK, boardGame)
}

// Method that Gets a boardgame based on a name
func (controller *Controller) GetAll(w http.ResponseWriter, r *http.Request) {

	boardgames, err := controller.repo.GetAll()
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	render.New().JSON(w, http.StatusOK, boardgames)
}

// Method that Gets a boardgame based on a name
func (controller *Controller) GetByName(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	boardgame, err := controller.repo.GetByName(name)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	render.New().JSON(w, http.StatusOK, boardgame)
}

// Method that Updates a boardgame based on an uuid and json input
func (controller *Controller) Update(w http.ResponseWriter, r *http.Request) {

	// Get Boardgame input from JSON input
	var input model.BoardGame
	err := utils.DecodeJSONBody(w, r, &input)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Get Boardgame by id
	id := utils.GetFieldFromURL(r, "id")
	boardgame, err := controller.repo.GetById(id)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Update Boardgame
	boardgame.UpdateBoardGame(input.GetName(), input.GetDealer(), input.GetPrice(), input.GetPlayerNumber())
	err = controller.repo.Update(boardgame)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	render.New().JSON(w, http.StatusOK, boardgame)
}

// Method that Deletes a boardgame based on an uuid
func (controller *Controller) Delete(w http.ResponseWriter, r *http.Request) {

	id := utils.GetFieldFromURL(r, "id")

	// Get Boardgame by name
	boardgame, err := controller.repo.GetById(id)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Delete by Id
	err = controller.repo.DeleteById(boardgame)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	render.New().JSON(w, http.StatusNoContent, id)
}
