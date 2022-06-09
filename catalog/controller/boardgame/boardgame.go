package boardgame

import (
	"log"
	"net/http"

	"catalog/model"
	"catalog/repository/boardgameRepo"
	"catalog/repository/categoryRepo"
	"catalog/repository/mechanismRepo"
	"catalog/repository/tagRepo"
	"catalog/utils"
)

// Declaring the repository interface in the controller package allows us to easily swap out the actual implementation, enforcing loose coupling.
type repository interface {
	Create(boardgame *model.Boardgame) error
	GetAll(sort, filterBody, filterValue string) ([]model.Boardgame, error)
	GetById(id string) (model.Boardgame, error)
	Update(boardgame model.Boardgame) error
	DeleteById(boardgame model.Boardgame) error
}

type tagRepository interface {
	Get(name string) (model.Tag, error)
}

type categoryRepository interface {
	Get(name string) (model.Category, error)
}

type mechanismRepository interface {
	Get(name string) (model.Mechanism, error)
}

// Controller contains the service, which contains database-related logic, as an injectable dependency, allowing us to decouple business logic from db logic.
type Controller struct {
	repo      repository
	tag       tagRepository
	category  categoryRepository
	mechanism mechanismRepository
}

// InitController initializes the boargame and the associations controller.
func InitController(boardGameRepo *boardgameRepo.BoardGameRepository, tagRepo *tagRepo.TagRepository, categoryRepo *categoryRepo.CategoryRepository, mechanismRepo *mechanismRepo.MechanismRepository) *Controller {
	return &Controller{
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
func (controller *Controller) Create(w http.ResponseWriter, r *http.Request) {

	var boardgame model.Boardgame
	err := utils.DecodeJSONBody(w, r, &boardgame)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	// Validate Boardgame input
	err = utils.ValidateStruct(&boardgame)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	// Check if Expansion -> Connect if needed
	id := utils.GetFieldFromURL(r, "id")
	err = controller.connectBoardgameToExpansion(id, &boardgame)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	// Check if Tags, Categories & Mechanisms exist
	err = controller.validateAssociations(w, r, boardgame)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	log.Println(boardgame)

	err = controller.repo.Create(&boardgame)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	log.Println(boardgame)

	utils.HTTPHandler(w, &boardgame, http.StatusOK, nil)
}

// Get Boardgames godoc
// @Summary 	Fetches all Boardgames
// @Tags 		boardgames
// @Produce 	json
// @Param 		filterBy query string  false  "Filter using field.value (For String partial find) OR field.operator.value"
// @Success 	200 {object} model.Boardgame
// @Router 		/boardgame [get]
func (controller *Controller) GetAll(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sortBy")
	sort, err := utils.GetSort(model.Boardgame{}, sortBy)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	filterBy := r.URL.Query().Get("filterBy")
	filterBody, filterValue, err := utils.GetFilters(model.Boardgame{}, filterBy)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	boardgames, err := controller.repo.GetAll(sort, filterBody, filterValue)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &boardgames, http.StatusOK, nil)
}

// Get Boardgame by id godoc
// @Summary 	Fetches a specific Boardgame using an id
// @Tags 		boardgames
// @Produce 	json
// @Param 		id path int true "The Boardgame unique id"
// @Success 	200 {object} model.Boardgame
// @Router 		/boardgame/{id} [get]
func (controller *Controller) Get(w http.ResponseWriter, r *http.Request) {

	id := utils.GetFieldFromURL(r, "id")

	boardgame, err := controller.repo.GetById(id)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &boardgame, http.StatusOK, nil)
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
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	// Validate Boardgame input
	err = utils.ValidateStruct(&input)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	// Check if Tags & Categories & Mechanisms exist
	err = controller.validateAssociations(w, r, input)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	// Get Boardgame by id
	id := utils.GetFieldFromURL(r, "id")
	boardgame, err := controller.repo.GetById(id)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	// Updates Boardgame
	boardgame.UpdateBoardgame(input)
	err = controller.repo.Update(boardgame)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &boardgame, http.StatusOK, nil)
}

// Delete Boardgame by id godoc
// @Summary 	Deletes a specific Boardgame via Id
// @Tags 		boardgames
// @Produce 	json
// @Param 		id path int true "The Boardgame id"
// @Success 	204
// @Router 		/boardgame/{id} [delete]
func (controller *Controller) Delete(w http.ResponseWriter, r *http.Request) {

	id := utils.GetFieldFromURL(r, "id")

	boardgame, err := controller.repo.GetById(id)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	// Delete by Id
	err = controller.repo.DeleteById(boardgame)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, id, http.StatusNoContent, nil)
}

// Function that checks if we are dealing with expansions and creates connection to boardgame parent if yes
func (controller *Controller) connectBoardgameToExpansion(id string, boardgame *model.Boardgame) error {
	if id != "" { // This is an expansion
		boardgameParent, err := controller.repo.GetById(id) // Get Parent BG
		if err != nil {
			return err
		}

		if boardgameParent.IsExpansion() {
			log.Println("Error -> An expansion cannot have other expansions")
			return utils.NewError(http.StatusUnprocessableEntity, " Error occurred while creating connection between Boardgames -> Expansion can't have expansions")
		}

		boardgame.SetBoardgameID(boardgameParent.GetId()) // Set the Parents Id in the expansion
	}
	return nil
}

// <<<<<<<<<<<< I DISLIKE THIS APPROACH, NOT MODULAR AND WILL WANT TO CHANGE >>>>>>>>>>>>>>>>>>

// Function that validates if tags and categories exist when boardgames are created
func (controller *Controller) validateAssociations(w http.ResponseWriter, r *http.Request, boardgame model.Boardgame) error {

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
