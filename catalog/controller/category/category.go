package category

import (
	"net/http"

	"catalog/model"
	"catalog/repository/categoryRepo"
	"catalog/utils"
)

type repository interface {
	Create(category model.Category) error
	GetAll(sort string) ([]model.Category, error)
	Get(name string) (model.Category, error)
	Delete(category model.Category) error
}

type Controller struct {
	repo repository
}

// InitController initializes the category controller.
func InitController(categoryRepo *categoryRepo.CategoryRepository) *Controller {
	return &Controller{
		repo: categoryRepo,
	}
}

// Create Category godoc
// @Summary 	Creates a Category using a name
// @Tags 		categories
// @Produce 	json
// @Param 		data body model.Category true "The Category name"
// @Success 	200 {object} model.Category
// @Router 		/category/{name} [post]
func (controller *Controller) Create(w http.ResponseWriter, r *http.Request) {

	var category model.Category
	err := utils.DecodeJSONBody(w, r, &category)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	// Validate Category input
	err = utils.ValidateStruct(&category)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	err = controller.repo.Create(category)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &category, http.StatusOK, nil)
}

// Get Categories godoc
// @Summary 	Fetches all Categories
// @Tags 		categories
// @Produce 	json
// @Success 	200 {object} model.Category
// @Router 		/category [get]
func (controller *Controller) GetAll(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sortBy")
	sort, err := utils.GetSort(model.Category{}, sortBy)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	categories, err := controller.repo.GetAll(sort)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &categories, http.StatusOK, nil)
}

// Get Category godoc
// @Summary 	Fetches a specific Category using a name
// @Tags 		categories
// @Produce 	json
// @Param 		name path string true "The Category name"
// @Success 	200 {object} model.Category
// @Router 		/category/{name} [get]
func (controller *Controller) Get(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	category, err := controller.repo.Get(name)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &category, http.StatusOK, nil)
}

// Delete Category godoc
// @Summary 	Deletes a specific Category
// @Tags 		categories
// @Produce 	json
// @Param 		name path string true "The Category name"
// @Success 	204
// @Router 		/category/{name} [delete]
func (controller *Controller) Delete(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	// Get category by name
	category, err := controller.repo.Get(name)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	// Delete by id
	err = controller.repo.Delete(category)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}
	utils.HTTPHandler(w, name, http.StatusNoContent, nil)
}
