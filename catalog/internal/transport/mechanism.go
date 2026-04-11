package transport

import (
	"context"
	"net/http"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/utils"
)

//go:generate mockgen -package transport -destination mechanism_mock.go . MechanismService

// MechanismService defines the interface for mechanism business logic.
type MechanismService interface {
	Create(ctx context.Context, mechanism *mechanism.Mechanism) error
	GetAll(ctx context.Context, sort string) ([]mechanism.Mechanism, error)
	Get(ctx context.Context, name string) (mechanism.Mechanism, error)
	Delete(ctx context.Context, name string) error
}

// MechanismController handles HTTP requests for mechanism operations.
type MechanismController struct {
	service MechanismService
}

// NewMechanismController initializes the mechanism controller.
func NewMechanismController(mechanismSvc MechanismService) *MechanismController {
	return &MechanismController{
		service: mechanismSvc,
	}
}

// Create Mechanism godoc
// @Summary 	Creates a Mechanism using a name
// @Tags 	mechanisms
// @Produce 	json
// @Param 		data body mechanism.Mechanism true "The Mechanism name"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} mechanism.Mechanism
// @Router 		/mechanism [post]
func (controller *MechanismController) Create(w http.ResponseWriter, r *http.Request) {
	// Deserialize Mechanism input
	var m = &mechanism.Mechanism{}
	if err := utils.DecodeJSONBody(w, r, m); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate Mechanism input
	if err := utils.ValidateStruct(m); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.service.Create(context.Background(), m); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, m); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// GetAll Mechanisms godoc
// @Summary 	Fetches all Mechanisms
// @Tags 	mechanisms
// @Produce 	json
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} mechanism.Mechanism
// @Router 		/mechanism [get]
func (controller *MechanismController) GetAll(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	sort, err := utils.GetSort(mechanism.Mechanism{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	mechanisms, err := controller.service.GetAll(context.Background(), sort)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, mechanisms); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// Get Mechanism godoc
// @Summary 	Fetches a specific Mechanism using a name
// @Tags 	mechanisms
// @Produce 	json
// @Param 		name path string true "The Mechanism name"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} mechanism.Mechanism
// @Router 		/mechanism/{name} [get]
func (controller *MechanismController) Get(w http.ResponseWriter, r *http.Request) {
	name := utils.GetFieldFromURL(r, "name")

	m, err := controller.service.Get(context.Background(), name)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, m); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// Delete Mechanism godoc
// @Summary 	Deletes a specific Mechanism
// @Tags 	mechanisms
// @Produce 	json
// @Param 		name path string true "The Mechanism name"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	204
// @Router 		/mechanism/{name} [delete]
func (controller *MechanismController) Delete(w http.ResponseWriter, r *http.Request) {
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
