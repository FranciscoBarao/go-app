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
	GetAll(ctx context.Context, opts ...listopt.Option) ([]contributor.Contributor, error)
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
// @Summary 	Lists contributors
// @Tags 		contributors
// @Produce 	json
// @Success 	200 {array} contributor.Contributor
// @Router 		/contributors [get]
func (c *ContributorController) GetAll(w http.ResponseWriter, r *http.Request) {
	var opts []listopt.Option

	sortBy := r.URL.Query().Get("sortBy")
	col, order, err := utils.GetSort(contributor.Contributor{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if col != "" {
		opts = append(opts, listopt.WithSort(col, order))
	}

	filterBy := r.URL.Query().Get("filterBy")
	fcol, fop, fval, err := utils.GetFilter(contributor.Contributor{}, filterBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if fcol != "" {
		opts = append(opts, listopt.WithFilter(fcol, fop, fval))
	}

	list, err := c.service.GetAll(r.Context(), opts...)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, list); err != nil {
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
