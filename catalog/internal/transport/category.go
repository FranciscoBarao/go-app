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
	GetAll(ctx context.Context, params listopt.Params) ([]category.Category, int, error)
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
// @Summary 	Lists Categories with pagination, sorting, and filtering
// @Tags 		categories
// @Produce 	json
// @Param 		page query int false "Page number (default 1)"
// @Param 		pageSize query int false "Items per page (default 10, max 100)"
// @Param 		sort query string false "Sort as field.order, e.g. name.asc"
// @Param 		filter query []string false "Filter as field.op.value, e.g. name.like.eco; repeatable, AND-combined"
// @Success 	200 {object} CategoryPage
// @Router 		/category [get]
func (controller *CategoryController) GetAll(w http.ResponseWriter, r *http.Request) {
	q, err := newQueryRequestFromURL(r.URL.Query())
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	controller.list(w, r, q)
}

// Query lists Categories with filtering, sorting, and pagination.
//
// It handles HTTP QUERY /api/category with a JSON body (pagination, sort,
// filters), for filters too complex to express as GET query parameters. Not
// represented in Swagger because OpenAPI 2.0 has no QUERY method; see GetAll for
// the same envelope.
func (controller *CategoryController) Query(w http.ResponseWriter, r *http.Request) {
	var q QueryRequest
	if err := utils.DecodeJSONBody(w, r, &q); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	controller.list(w, r, q)
}

// list resolves a validated list request and writes the paginated envelope,
// shared by the GET and QUERY entry points.
func (controller *CategoryController) list(w http.ResponseWriter, r *http.Request, q QueryRequest) {
	opts, err := q.toOptions(category.QuerySchema)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	p := listopt.Apply(opts...)
	categories, total, err := controller.service.GetAll(r.Context(), p)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, newPaginatedResponse(categories, total, p.Pagination)); err != nil {
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
