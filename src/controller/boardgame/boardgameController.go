package boardgame

import (
	"errors"
	"log"
	"net/http"

	"github.com/unrolled/render"

	"go-app/model"
	"go-app/repository/boardgameRepo"
	"go-app/repository/tagRepo"
	"go-app/utils"
)

// Declaring the repository interface in the controller package allows us to easily swap out the actual implementation, enforcing loose coupling.
type repository interface {
	Create(boardgame model.Boardgame) error
	GetAll() ([]model.Boardgame, error)
	GetByName(name string) (model.Boardgame, error)
	GetById(id string) (model.Boardgame, error)
	Update(boardgame model.Boardgame) error
	DeleteById(boardgame model.Boardgame) error
}

type tagRepository interface {
	GetByName(name string) (model.Tag, error)
}

// Controller contains the service, which contains database-related logic, as an injectable dependency, allowing us to decouple business logic from db logic.
type Controller struct {
	repo repository
	tag  tagRepository
}

// InitController initializes the boargame controller.
func InitController(boardGameRepo *boardgameRepo.BoardGameRepository, tagRepo *tagRepo.TagRepository) *Controller {
	return &Controller{
		repo: boardGameRepo,
		tag:  tagRepo,
	}
}

// Method that Creates a boardgame based on json input
func (controller *Controller) Create(w http.ResponseWriter, r *http.Request) {

	var boardgame model.Boardgame
	err := utils.DecodeJSONBody(w, r, &boardgame)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Check if Tags exist
	err = controller.validateTags(w, r, boardgame)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = controller.repo.Create(boardgame)
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
	var input model.Boardgame
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

	// Check if Tags exist
	err = controller.validateTags(w, r, input)
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

	// Updates Boardgame
	boardgame.UpdateBoardgame(input.GetName(), input.GetDealer(), input.GetPrice(), input.GetPlayerNumber(), input.GetTags())
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

func (controller *Controller) validateTags(w http.ResponseWriter, r *http.Request, boardgame model.Boardgame) error {

	// Boardgame can contain Tags ->  We omit them which means that if they don't previously exist, the db returns an error -> Check if they exist before hand
	if boardgame.IsTags() {
		for _, tempTag := range boardgame.GetTags() {

			tag, err := controller.tag.GetByName(tempTag.GetName()) // Get tag by name
			if err != nil {                                         // That tag does not exist -> Return Error
				return err
			}
			log.Println("TAG: " + tag.GetName())
		}
	}
	return nil
}
