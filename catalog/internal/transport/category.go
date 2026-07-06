package transport

import (
	"context"
	"net/http"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/utils"
)

//go:generate mockgen -package transport -destination category_mock.go . CategoryService

// CategoryService defines the interface for category business logic.
type CategoryService interface {
	Create(ctx context.Context, req *category.CreateCategoryRequest) (category.Category, error)
	GetAll(ctx context.Context, opts ...listopt.Option) ([]category.Category, error)
	Get(ctx context.Context, slug string) (category.Category, error)
	Delete(ctx context.Context, slug string, hard bool) error
}

// CategoryController handles HTTP requests for category operations.
type CategoryController struct {
	service CategoryService
}

// NewCategoryController initializes the category controller.
func NewCategoryController(categorySvc CategoryService) *CategoryController {
	return &CategoryController{service: categorySvc}
}

// Create Category godoc
// @Summary 	Creates a Category
// @Tags 		categories
// @Produce 	json
// @Param 		data body category.CreateCategoryRequest true "The Category"
// @Success 	200 {object} category.Category
// @Router 		/category [post]
func (controller *CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	var req category.CreateCategoryRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := utils.ValidateStruct(&req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	c, err := controller.service.Create(r.Context(), &req)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, c); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// GetAll Categories godoc
// @Summary 	Fetches all Categories
// @Tags 		categories
// @Produce 	json
// @Success 	200 {object} category.Category
// @Router 		/category [get]
func (controller *CategoryController) GetAll(w http.ResponseWriter, r *http.Request) {
	var opts []listopt.Option

	sortBy := r.URL.Query().Get("sortBy")
	col, order, err := utils.GetSort(category.Category{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if col != "" {
		opts = append(opts, listopt.WithSort(col, order))
	}

	filterBy := r.URL.Query().Get("filterBy")
	fcol, fop, fval, err := utils.GetFilter(category.Category{}, filterBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if fcol != "" {
		opts = append(opts, listopt.WithFilter(fcol, fop, fval))
	}

	categories, err := controller.service.GetAll(r.Context(), opts...)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, categories); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// Get Category godoc
// @Summary 	Fetches a Category by slug
// @Tags 		categories
// @Produce 	json
// @Param 		slug path string true "Category slug"
// @Success 	200 {object} category.Category
// @Router 		/category/{slug} [get]
func (controller *CategoryController) Get(w http.ResponseWriter, r *http.Request) {
	slugStr := utils.GetFieldFromURL(r, "slug")
	c, err := controller.service.Get(r.Context(), slugStr)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, c); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// Delete Category godoc
// @Summary 	Deletes a Category
// @Tags 		categories
// @Produce 	json
// @Param 		slug path string true "Category slug"
// @Param 		hard query bool false "Hard delete"
// @Success 	204
// @Router 		/category/{slug} [delete]
func (controller *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	slugStr := utils.GetFieldFromURL(r, "slug")
	hard := r.URL.Query().Get("hard") == "true"

	if err := controller.service.Delete(r.Context(), slugStr, hard); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
