package controllers

import (
	"net/http"

	"catalog/middleware"
	"catalog/model"
	"catalog/services"
	"catalog/utils"

	"github.com/unrolled/render"
)

// Declaring the repository interface in the controller package allows us to easily swap out the actual implementation, enforcing loose coupling.
type boardgameService interface {
	Create(boardgame *model.Boardgame, id string) error
	GetAll(sort, filterBody, filterValue string) ([]model.Boardgame, error)
	GetById(id string) (model.Boardgame, error)
	Update(boardgame model.Boardgame, id string) error
	DeleteById(id string) error
}

// Controller contains the service, which contains database-related logic, as an injectable dependency, allowing us to decouple business logic from db logic.
type BoardgameController struct {
	service boardgameService
}

// InitController initializes the boargame and the associations controller.
func InitBoardgameController(boardGameSvc *services.BoardgameService) *BoardgameController {
	return &BoardgameController{
		service: boardGameSvc,
	}
}

// Create Boardgame godoc
// @Summary 	Creates a Boardgame based on a json body
// @Tags 		boardgames
// @Produce 	json
// @Param 		data body model.Boardgame true "The input Boardgame struct"
// @Param 		id path int false "The Boardgame id indicating this is an Expansion"
// @Success 	200 {object} model.Boardgame
// @Router 		/boardgame [post]
func (controller *BoardgameController) Create(w http.ResponseWriter, r *http.Request) {

	// Deserialize Boardgame input
	var boardgame model.Boardgame
	if err := utils.DecodeJSONBody(w, r, &boardgame); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate Boardgame input
	if err := utils.ValidateStruct(&boardgame); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Get Id from url - If its an expansion
	id := utils.GetFieldFromURL(r, "id")

	if err := controller.service.Create(&boardgame, id); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, boardgame)
}

// Get Boardgames godoc
// @Summary 	Fetches all Boardgames
// @Tags 		boardgames
// @Produce 	json
// @Param 		filterBy query string  false  "Filter using field.value (For String partial find) OR field.operator.value"
// @Success 	200 {object} model.Boardgame
// @Router 		/boardgame [get]
func (controller *BoardgameController) GetAll(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sortBy")
	sort, err := utils.GetSort(model.Boardgame{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	filterBy := r.URL.Query().Get("filterBy")
	filterBody, filterValue, err := utils.GetFilters(model.Boardgame{}, filterBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	boardgames, err := controller.service.GetAll(sort, filterBody, filterValue)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, boardgames)

}

// Get Boardgame by id godoc
// @Summary 	Fetches a specific Boardgame using an id
// @Tags 		boardgames
// @Produce 	json
// @Param 		id path int true "The Boardgame unique id"
// @Success 	200 {object} model.Boardgame
// @Router 		/boardgame/{id} [get]
func (controller *BoardgameController) Get(w http.ResponseWriter, r *http.Request) {

	id := utils.GetFieldFromURL(r, "id")

	boardgame, err := controller.service.GetById(id)
	if err != nil {
		middleware.ErrorHandler(w, err)
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
func (controller *BoardgameController) Update(w http.ResponseWriter, r *http.Request) {

	// Deserialize Boardgame input
	var input model.Boardgame
	if err := utils.DecodeJSONBody(w, r, &input); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate Boardgame input
	if err := utils.ValidateStruct(&input); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	id := utils.GetFieldFromURL(r, "id")

	// Updates Boardgame
	if err := controller.service.Update(input, id); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, input)
}

// Delete Boardgame by id godoc
// @Summary 	Deletes a specific Boardgame via Id
// @Tags 		boardgames
// @Produce 	json
// @Param 		id path int true "The Boardgame id"
// @Success 	204
// @Router 		/boardgame/{id} [delete]
func (controller *BoardgameController) Delete(w http.ResponseWriter, r *http.Request) {

	id := utils.GetFieldFromURL(r, "id")

	// Delete by Id
	if err := controller.service.DeleteById(id); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusNoContent, id)
}
