package transport

import (
	"context"
	"net/http"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/query"
	"github.com/FranciscoBarao/catalog/internal/utils"
)

//go:generate mockgen -package transport -destination category_mock.go . CategoryService

// CategoryService defines the interface for category business logic.
type CategoryService interface {
	Create(ctx context.Context, c *category.Category) error
	GetAll(ctx context.Context, opts ...query.Option) ([]category.Category, error)
	Get(ctx context.Context, name string) (category.Category, error)
	Delete(ctx context.Context, name string) error
}

// CategoryController handles HTTP requests for category operations.
type CategoryController struct {
	service CategoryService
}

// NewCategoryController initializes the category controller.
func NewCategoryController(categorySvc CategoryService) *CategoryController {
	return &CategoryController{
		service: categorySvc,
	}
}

// Create Category godoc
// @Summary 	Creates a Category using a name
// @Tags 		categories
// @Produce 	json
// @Param 		data body category.Category true "The Category name"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} category.Category
// @Router 		/category [post]
func (controller *CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	// Deserialize Category input
	var c = &category.Category{}
	if err := utils.DecodeJSONBody(w, r, c); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate Category input
	if err := utils.ValidateStruct(c); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.service.Create(context.Background(), c); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, c); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// GetAll Categories godoc
// @Summary 	Fetches all Categories
// @Tags 		categories
// @Produce 	json
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} category.Category
// @Router 		/category [get]
func (controller *CategoryController) GetAll(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	col, order, err := utils.GetSort(category.Category{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	var opts []query.Option
	if col != "" {
		opts = append(opts, query.WithSort(col, order))
	}

	categories, err := controller.service.GetAll(context.Background(), opts...)
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
// @Success 	200 {object} category.Category
// @Router 		/category/{name} [get]
func (controller *CategoryController) Get(w http.ResponseWriter, r *http.Request) {
	name := utils.GetFieldFromURL(r, "name")
	c, err := controller.service.Get(context.Background(), name)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, c); err != nil {
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
	if err := controller.service.Delete(context.Background(), name); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusNoContent, name); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}
