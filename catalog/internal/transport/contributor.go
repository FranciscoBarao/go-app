package transport

import (
	"context"
	"net/http"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/utils"
)

//go:generate mockgen -package transport -destination contributor_mock.go . ContributorService

// ContributorService defines contributor business logic.
type ContributorService interface {
	Create(ctx context.Context, req *contributor.CreateContributorRequest) (contributor.Contributor, error)
	Update(ctx context.Context, req *contributor.UpdateContributorRequest, slug string) (contributor.Contributor, error)
	GetAll(ctx context.Context, opts ...listopt.Option) ([]contributor.Contributor, int, error)
	Get(ctx context.Context, slug string) (contributor.Contributor, error)
	Delete(ctx context.Context, slug string, hard bool) error
}

// ContributorController handles HTTP requests for contributors.
type ContributorController struct {
	service ContributorService
}

// NewContributorController creates a contributor controller.
func NewContributorController(svc ContributorService) *ContributorController {
	return &ContributorController{service: svc}
}

// Create Contributor godoc
// @Summary 	Creates a contributor
// @Tags 		contributors
// @Produce 	json
// @Param 		data body contributor.CreateContributorRequest true "Contributor"
// @Success 	200 {object} contributor.Contributor
// @Router 		/contributors [post]
func (c *ContributorController) Create(w http.ResponseWriter, r *http.Request) {
	var req contributor.CreateContributorRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := utils.ValidateStruct(&req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	contrib, err := c.service.Create(r.Context(), &req)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, contrib); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// GetAll Contributors godoc
// @Summary 	Lists Contributors with pagination, sorting, and filtering
// @Tags 		contributors
// @Produce 	json
// @Param 		page query int false "Page number (default 1)"
// @Param 		pageSize query int false "Items per page (default 10, max 100)"
// @Param 		sort query string false "Sort as field.order, e.g. name.asc"
// @Param 		filter query []string false "Filter as field.op.value, e.g. name.like.knizia; repeatable, AND-combined"
// @Success 	200 {object} ContributorPage
// @Router 		/contributors [get]
func (c *ContributorController) GetAll(w http.ResponseWriter, r *http.Request) {
	q, err := newQueryRequestFromURL(r.URL.Query())
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	c.list(w, r, q)
}

// Query lists contributors with filtering, sorting, and pagination.
//
// It handles HTTP QUERY /api/contributors with a JSON body (pagination, sort,
// filters), for filters too complex to express as GET query parameters. Not
// represented in Swagger because OpenAPI 2.0 has no QUERY method; see GetAll for
// the same envelope.
func (c *ContributorController) Query(w http.ResponseWriter, r *http.Request) {
	var q QueryRequest
	if err := utils.DecodeJSONBody(w, r, &q); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	c.list(w, r, q)
}

// list resolves a validated list request and writes the paginated envelope,
// shared by the GET and QUERY entry points.
func (c *ContributorController) list(w http.ResponseWriter, r *http.Request, q QueryRequest) {
	opts, err := q.toOptions(contributor.Contributor{})
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	p := listopt.Apply(opts...)
	list, total, err := c.service.GetAll(r.Context(), opts...)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, newPaginatedResponse(list, total, p.Pagination)); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// Get Contributor godoc
// @Summary 	Fetches a contributor by slug
// @Tags 		contributors
// @Produce 	json
// @Param 		slug path string true "Contributor slug"
// @Success 	200 {object} contributor.Contributor
// @Router 		/contributors/{slug} [get]
func (c *ContributorController) Get(w http.ResponseWriter, r *http.Request) {
	slugStr := utils.GetFieldFromURL(r, "slug")
	contrib, err := c.service.Get(r.Context(), slugStr)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, contrib); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// Update Contributor godoc
// @Summary 	Updates a contributor
// @Tags 		contributors
// @Produce 	json
// @Param 		slug path string true "Contributor slug"
// @Param 		data body contributor.UpdateContributorRequest true "Contributor fields"
// @Success 	200 {object} contributor.Contributor
// @Router 		/contributors/{slug} [patch]
func (c *ContributorController) Update(w http.ResponseWriter, r *http.Request) {
	slugStr := utils.GetFieldFromURL(r, "slug")
	var req contributor.UpdateContributorRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := utils.ValidateStruct(&req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	contrib, err := c.service.Update(r.Context(), &req, slugStr)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, contrib); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// Delete Contributor godoc
// @Summary 	Deletes a contributor
// @Tags 		contributors
// @Produce 	json
// @Param 		slug path string true "Contributor slug"
// @Param 		hard query bool false "Hard delete"
// @Success 	204
// @Router 		/contributors/{slug} [delete]
func (c *ContributorController) Delete(w http.ResponseWriter, r *http.Request) {
	slugStr := utils.GetFieldFromURL(r, "slug")
	hard := r.URL.Query().Get("hard") == "true"

	if err := c.service.Delete(r.Context(), slugStr, hard); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
