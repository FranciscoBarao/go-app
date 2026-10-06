package transport

import (
	"context"
	"errors"
	"net/http"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/utils"
)

//go:generate mockgen -package transport -destination mechanism_mock.go . MechanismService

// MechanismService defines the interface for mechanism business logic.
type MechanismService interface {
	Create(ctx context.Context, req *mechanism.CreateMechanismRequest) (mechanism.Mechanism, error)
	GetAll(ctx context.Context, query listopt.Query) ([]mechanism.Mechanism, int, error)
	Get(ctx context.Context, slug string) (mechanism.Mechanism, error)
	Delete(ctx context.Context, slug string, hard bool) error
}

// MechanismController handles HTTP requests for mechanism operations.
type MechanismController struct {
	service MechanismService
}

// NewMechanismController initializes the mechanism controller.
func NewMechanismController(mechanismSvc MechanismService) *MechanismController {
	return &MechanismController{service: mechanismSvc}
}

// Create Mechanism godoc
// @Summary 	Creates a Mechanism
// @Tags 		mechanisms
// @Produce 	json
// @Param 		data body mechanism.CreateMechanismRequest true "The Mechanism"
// @Success 	200 {object} mechanism.Mechanism
// @Router 		/mechanism [post]
func (controller *MechanismController) Create(w http.ResponseWriter, r *http.Request) {
	var req mechanism.CreateMechanismRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		switch {
		case errors.Is(err, utils.ErrBadRequest):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, utils.ErrRequestTooLarge):
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("decode create mechanism")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := utils.ValidateStruct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m, err := controller.service.Create(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, mechanism.ErrValidation):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, mechanism.ErrAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("create mechanism")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := render.New().JSON(w, http.StatusOK, m); err != nil {
		logging.FromCtx(r.Context()).Error().Err(err).Msg("encode mechanism")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetAll Mechanisms godoc
// @Summary 	Lists Mechanisms with pagination, sorting, and filtering
// @Tags 		mechanisms
// @Produce 	json
// @Param 		page query int false "Page number (default 1)"
// @Param 		pageSize query int false "Items per page (default 10, max 100)"
// @Param 		sort query string false "Sort as field.order, e.g. name.asc"
// @Param 		filter query []string false "Filter as field.op.value, e.g. name.like.trad; repeatable, AND-combined"
// @Success 	200 {object} MechanismPage
// @Router 		/mechanism [get]
func (controller *MechanismController) GetAll(w http.ResponseWriter, r *http.Request) {
	q, err := newQueryRequestFromURL(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	controller.list(w, r, q)
}

// Query lists Mechanisms with filtering, sorting, and pagination.
//
// It handles HTTP QUERY /api/mechanism with a JSON body (pagination, sort,
// filters), for filters too complex to express as GET query parameters. Not
// represented in Swagger because OpenAPI 2.0 has no QUERY method; see GetAll for
// the same envelope.
func (controller *MechanismController) Query(w http.ResponseWriter, r *http.Request) {
	var q QueryRequest
	if err := utils.DecodeJSONBody(w, r, &q); err != nil {
		switch {
		case errors.Is(err, utils.ErrBadRequest):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, utils.ErrRequestTooLarge):
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("decode mechanism query")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	controller.list(w, r, q)
}

// list resolves a validated list request and writes the paginated envelope,
// shared by the GET and QUERY entry points.
func (controller *MechanismController) list(w http.ResponseWriter, r *http.Request, q QueryRequest) {
	query, err := q.toQuery(mechanism.QueryAllowlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	mechanisms, total, err := controller.service.GetAll(r.Context(), query)
	if err != nil {
		logging.FromCtx(r.Context()).Error().Err(err).Msg("list mechanisms")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, newPaginatedResponse(mechanisms, total, query.Pagination)); err != nil {
		logging.FromCtx(r.Context()).Error().Err(err).Msg("encode mechanisms")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Get Mechanism godoc
// @Summary 	Fetches a Mechanism by slug
// @Tags 		mechanisms
// @Produce 	json
// @Param 		slug path string true "Mechanism slug"
// @Success 	200 {object} mechanism.Mechanism
// @Router 		/mechanism/{slug} [get]
func (controller *MechanismController) Get(w http.ResponseWriter, r *http.Request) {
	slugStr := utils.GetFieldFromURL(r, "slug")
	m, err := controller.service.Get(r.Context(), slugStr)
	if err != nil {
		switch {
		case errors.Is(err, mechanism.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("get mechanism")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := render.New().JSON(w, http.StatusOK, m); err != nil {
		logging.FromCtx(r.Context()).Error().Err(err).Msg("encode mechanism")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Delete Mechanism godoc
// @Summary 	Deletes a Mechanism
// @Tags 		mechanisms
// @Produce 	json
// @Param 		slug path string true "Mechanism slug"
// @Param 		hard query bool false "Hard delete"
// @Success 	204
// @Router 		/mechanism/{slug} [delete]
func (controller *MechanismController) Delete(w http.ResponseWriter, r *http.Request) {
	slugStr := utils.GetFieldFromURL(r, "slug")
	hard := r.URL.Query().Get("hard") == "true"

	if err := controller.service.Delete(r.Context(), slugStr, hard); err != nil {
		switch {
		case errors.Is(err, mechanism.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, mechanism.ErrStillReferenced):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("delete mechanism")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
