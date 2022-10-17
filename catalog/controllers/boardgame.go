package controllers

import (
	"log"
	"net/http"

	"catalog/middleware"
	"catalog/model"
	"catalog/repositories"
	"catalog/utils"

	"github.com/unrolled/render"
)

// Declaring the repository interface in the controller package allows us to easily swap out the actual implementation, enforcing loose coupling.
type boardgameRepository interface {
	Create(boardgame *model.Boardgame) error
	GetAll(sort, filterBody, filterValue string) ([]model.Boardgame, error)
	GetById(id string) (model.Boardgame, error)
	Update(boardgame model.Boardgame) error
	DeleteById(boardgame model.Boardgame) error
}

// Controller contains the service, which contains database-related logic, as an injectable dependency, allowing us to decouple business logic from db logic.
type BoardgameController struct {
	repo      boardgameRepository
	tag       tagRepository
	category  categoryRepository
	mechanism mechanismRepository
}

// InitController initializes the boargame and the associations controller.
func InitBoardgameController(boardGameRepo *repositories.BoardgameRepository, tagRepo *repositories.TagRepository, categoryRepo *repositories.CategoryRepository, mechanismRepo *repositories.MechanismRepository) *BoardgameController {
	return &BoardgameController{
		repo:      boardGameRepo,
		tag:       tagRepo,
		category:  categoryRepo,
		mechanism: mechanismRepo,
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

	// Check if Expansion -> Connect if needed
	id := utils.GetFieldFromURL(r, "id")
	if err := controller.connectBoardgameToExpansion(id, &boardgame); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Check if Tags, Categories & Mechanisms exist
	if err := controller.validateAssociations(w, r, boardgame); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.repo.Create(&boardgame); err != nil {
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

	boardgames, err := controller.repo.GetAll(sort, filterBody, filterValue)
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

	boardgame, err := controller.repo.GetById(id)
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

	// Check if Tags & Categories & Mechanisms exist
	if err := controller.validateAssociations(w, r, input); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Get Boardgame by id
	id := utils.GetFieldFromURL(r, "id")
	boardgame, err := controller.repo.GetById(id)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Updates Boardgame
	boardgame.UpdateBoardgame(input)
	if err := controller.repo.Update(boardgame); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, boardgame)
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

	boardgame, err := controller.repo.GetById(id)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Delete by Id
	if err := controller.repo.DeleteById(boardgame); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusNoContent, id)
}

// Function that checks if we are dealing with expansions and creates connection to boardgame parent if yes
func (controller *BoardgameController) connectBoardgameToExpansion(id string, boardgame *model.Boardgame) error {
	if id != "" { // This is an expansion
		boardgameParent, err := controller.repo.GetById(id) // Get Parent BG
		if err != nil {
			return err
		}

		if boardgameParent.IsExpansion() {
			log.Println("Error -> An expansion cannot have other expansions")
			return middleware.NewError(http.StatusConflict, "Expansion can't have expansions")
		}

		boardgame.SetBoardgameID(boardgameParent.GetId()) // Set the Parents Id in the expansion
	}
	return nil
}

// <<<<<<<<<<<< I DISLIKE THIS APPROACH, NOT MODULAR AND WILL WANT TO CHANGE >>>>>>>>>>>>>>>>>>

// Function that validates if tags and categories exist when boardgames are created
func (controller *BoardgameController) validateAssociations(w http.ResponseWriter, r *http.Request, boardgame model.Boardgame) error {

	// Boardgame can contain Associations like Tags or Categories ->  We omit them which means that if they don't previously exist, the db returns an error -> Check if they exist before hand
	if boardgame.IsTags() {
		for _, tempTag := range boardgame.GetTags() {

			_, err := controller.tag.Get(tempTag.GetName()) // Get tag by name
			if err != nil {                                 // That tag does not exist -> Return Error
				return err
			}
		}
	}

	if boardgame.IsCategories() {
		for _, tempCategory := range boardgame.GetCategories() {

			_, err := controller.category.Get(tempCategory.GetName()) // Get category by name
			if err != nil {                                           // That category does not exist -> Return Error
				return err
			}
		}
	}

	if boardgame.IsMechanisms() {
		for _, tempMechanism := range boardgame.GetMechanisms() {

			_, err := controller.mechanism.Get(tempMechanism.GetName()) // Get mechanism by name
			if err != nil {                                             // That mechanism does not exist -> Return Error
				return err
			}
		}
	}
	return nil
}
