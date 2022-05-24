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
	Get(name string) (model.Tag, error)
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

// Create Boardgame godoc
// @Summary 	Creates a Boardgame based on a json body
// @Tags 		boardgames
// @Produce 	json
// @Param 		data body model.Boardgame true "The input Boardgame struct"
// @Success 	200 {object} model.Boardgame
// @Router 		/boardgame [post]
func (controller *Controller) Create(w http.ResponseWriter, r *http.Request) {
	// Method that Creates a boardgame based on json input

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

// Get Boardgames godoc
// @Summary 	Fetches all Boardgames
// @Tags 		boardgames
// @Produce 	json
// @Success 	200 {object} model.Boardgame
// @Router 		/boardgame [get]
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

// Get Boardgame by name godoc
// @Summary 	Fetches a specific Boardgame using a name
// @Tags 		boardgames
// @Produce 	json
// @Param 		name path string true "The Boardgame name"
// @Success 	200 {object} model.Boardgame
// @Router 		/boardgame/{name} [get]
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

// Update Boardgame by id godoc
// @Summary 	Updates a specific Boardgame via Id
// @Tags 		boardgames
// @Produce 	json
// @Param 		id path int true "The Boardgame id"
// @Param 		data body model.Boardgame true "The Boardgame struct to be updated into"
// @Success 	200 {object} model.Boardgame
// @Router 		/boardgame/{id} [patch]
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
	boardgame.UpdateBoardgame(input.GetName(), input.GetPublisher(), input.GetPrice(), input.GetPlayerNumber(), input.GetTags())
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

// Delete Boardgame by id godoc
// @Summary 	Deletes a specific Boardgame via Id
// @Tags 		boardgames
// @Produce 	json
// @Param 		id path int true "The Boardgame id"
// @Success 	200
// @Router 		/boardgame/{id} [delete]
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

			tag, err := controller.tag.Get(tempTag.GetName()) // Get tag by name
			if err != nil {                                   // That tag does not exist -> Return Error
				return err
			}
			log.Println("TAG: " + tag.GetName())
		}
	}
	return nil
}
