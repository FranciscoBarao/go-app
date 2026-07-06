package transport

import (
	"context"
	"net/http"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/utils"
)

//go:generate mockgen -package transport -destination mechanism_mock.go . MechanismService

// MechanismService defines the interface for mechanism business logic.
type MechanismService interface {
	Create(ctx context.Context, req *mechanism.CreateMechanismRequest) (mechanism.Mechanism, error)
	GetAll(ctx context.Context, opts ...listopt.Option) ([]mechanism.Mechanism, error)
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
		middleware.ErrorHandler(w, err)
		return
	}
	if err := utils.ValidateStruct(&req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	m, err := controller.service.Create(r.Context(), &req)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, m); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// GetAll Mechanisms godoc
// @Summary 	Fetches all Mechanisms
// @Tags 		mechanisms
// @Produce 	json
// @Success 	200 {object} mechanism.Mechanism
// @Router 		/mechanism [get]
func (controller *MechanismController) GetAll(w http.ResponseWriter, r *http.Request) {
	var opts []listopt.Option

	sortBy := r.URL.Query().Get("sortBy")
	col, order, err := utils.GetSort(mechanism.Mechanism{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if col != "" {
		opts = append(opts, listopt.WithSort(col, order))
	}

	filterBy := r.URL.Query().Get("filterBy")
	fcol, fop, fval, err := utils.GetFilter(mechanism.Mechanism{}, filterBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if fcol != "" {
		opts = append(opts, listopt.WithFilter(fcol, fop, fval))
	}

	mechanisms, err := controller.service.GetAll(r.Context(), opts...)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, mechanisms); err != nil {
		middleware.ErrorHandler(w, err)
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
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, m); err != nil {
		middleware.ErrorHandler(w, err)
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
		middleware.ErrorHandler(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
