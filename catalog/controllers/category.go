package controllers

import (
	"net/http"

	"catalog/middleware"
	"catalog/model"
	"catalog/repositories"
	"catalog/utils"

	"github.com/unrolled/render"
)

type categoryRepository interface {
	Create(category model.Category) error
	GetAll(sort string) ([]model.Category, error)
	Get(name string) (model.Category, error)
	Delete(category model.Category) error
}

type CategoryController struct {
	repo categoryRepository
}

// InitController initializes the category controller.
func InitCategoryController(categoryRepo *repositories.CategoryRepository) *CategoryController {
	return &CategoryController{
		repo: categoryRepo,
	}
}

// Create Category godoc
// @Summary 	Creates a Category using a name
// @Tags 		categories
// @Produce 	json
// @Param 		data body model.Category true "The Category name"
// @Success 	200 {object} model.Category
// @Router 		/category [post]
func (controller *CategoryController) Create(w http.ResponseWriter, r *http.Request) {

	// Deserialize Category input
	var category model.Category
	if err := utils.DecodeJSONBody(w, r, &category); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate Category input
	if err := utils.ValidateStruct(&category); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.repo.Create(category); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, category)
}

// Get Categories godoc
// @Summary 	Fetches all Categories
// @Tags 		categories
// @Produce 	json
// @Success 	200 {object} model.Category
// @Router 		/category [get]
func (controller *CategoryController) GetAll(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sortBy")
	sort, err := utils.GetSort(model.Category{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	categories, err := controller.repo.GetAll(sort)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	render.New().JSON(w, http.StatusOK, categories)
}

// Get Category godoc
// @Summary 	Fetches a specific Category using a name
// @Tags 		categories
// @Produce 	json
// @Param 		name path string true "The Category name"
// @Success 	200 {object} model.Category
// @Router 		/category/{name} [get]
func (controller *CategoryController) Get(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	category, err := controller.repo.Get(name)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	render.New().JSON(w, http.StatusOK, category)
}

// Delete Category godoc
// @Summary 	Deletes a specific Category
// @Tags 		categories
// @Produce 	json
// @Param 		name path string true "The Category name"
// @Success 	204
// @Router 		/category/{name} [delete]
func (controller *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	// Get category by name
	category, err := controller.repo.Get(name)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Delete by id
	if err := controller.repo.Delete(category); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusNoContent, name)
}
