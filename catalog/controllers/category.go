package controllers

import (
	"net/http"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/model"
	"github.com/FranciscoBarao/catalog/services"
	"github.com/FranciscoBarao/catalog/utils"
)

type categoryService interface {
	Create(category *model.Category) error
	GetAll(sort string) ([]model.Category, error)
	Get(name string) (model.Category, error)
	Delete(name string) error
}

type CategoryController struct {
	service categoryService
}

// InitController initializes the category controller.
func InitCategoryController(categorySvc *services.CategoryService) *CategoryController {
	return &CategoryController{
		service: categorySvc,
	}
}

// Create Category godoc
// @Summary 	Creates a Category using a name
// @Tags 		categories
// @Produce 	json
// @Param 		data body model.Category true "The Category name"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} model.Category
// @Router 		/category [post]
func (controller *CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	// Deserialize Category input
	var category = &model.Category{}
	if err := utils.DecodeJSONBody(w, r, category); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate Category input
	if err := utils.ValidateStruct(category); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.service.Create(category); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, category); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// Get Categories godoc
// @Summary 	Fetches all Categories
// @Tags 		categories
// @Produce 	json
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} model.Category
// @Router 		/category [get]
func (controller *CategoryController) GetAll(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	sort, err := utils.GetSort(model.Category{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	categories, err := controller.service.GetAll(sort)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, categories); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// Get Category godoc
// @Summary 	Fetches a specific Category using a name
// @Tags 		categories
// @Produce 	json
// @Param 		name path string true "The Category name"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} model.Category
// @Router 		/category/{name} [get]
func (controller *CategoryController) Get(w http.ResponseWriter, r *http.Request) {
	name := utils.GetFieldFromURL(r, "name")
	category, err := controller.service.Get(name)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, category); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// Delete Category godoc
// @Summary 	Deletes a specific Category
// @Tags 		categories
// @Produce 	json
// @Param 		name path string true "The Category name"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	204
// @Router 		/category/{name} [delete]
func (controller *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	name := utils.GetFieldFromURL(r, "name")

	// Delete by id
	if err := controller.service.Delete(name); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusNoContent, name); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}
