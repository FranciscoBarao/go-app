package controllers

import (
	"net/http"

	"catalog/middleware"
	"catalog/model"
	"catalog/services"
	"catalog/utils"

	"github.com/unrolled/render"
)

type mechanismService interface {
	Create(mechanism *model.Mechanism) error
	GetAll(sort string) ([]model.Mechanism, error)
	Get(name string) (model.Mechanism, error)
	Delete(name string) error
}

type MechanismController struct {
	service mechanismService
}

// InitController initializes the mechanism controller.
func InitMechanismController(mechanismSvc *services.MechanismService) *MechanismController {
	return &MechanismController{
		service: mechanismSvc,
	}
}

// Create Mechanism godoc
// @Summary 	Creates a Mechanism using a name
// @Mechanisms 	mechanisms
// @Produce 	json
// @Param 		data body model.Mechanism true "The Mechanism name"
// @Success 	200 {object} model.Mechanism
// @Router 		/mechanism [post]
func (controller *MechanismController) Create(w http.ResponseWriter, r *http.Request) {

	// Deserialize Mechanism input
	var mechanism *model.Mechanism
	if err := utils.DecodeJSONBody(w, r, mechanism); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate Mechanism input
	if err := utils.ValidateStruct(mechanism); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.service.Create(mechanism); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, mechanism)
}

// Get Mechanisms godoc
// @Summary 	Fetches all Mechanisms
// @Mechanisms 	mechanisms
// @Produce 	json
// @Success 	200 {object} model.Mechanism
// @Router 		/mechanism [get]
func (controller *MechanismController) GetAll(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sortBy")
	sort, err := utils.GetSort(model.Mechanism{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	mechanisms, err := controller.service.GetAll(sort)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	render.New().JSON(w, http.StatusOK, mechanisms)
}

// Get Mechanism godoc
// @Summary 	Fetches a specific Mechanism using a name
// @Mechanisms 	mechanisms
// @Produce 	json
// @Param 		name path string true "The Mechanism name"
// @Success 	200 {object} model.Mechanism
// @Router 		/mechanism/{name} [get]
func (controller *MechanismController) Get(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	mechanism, err := controller.service.Get(name)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, mechanism)
}

// Delete Mechanism godoc
// @Summary 	Deletes a specific Mechanism
// @Mechanisms 	mechanisms
// @Produce 	json
// @Param 		name path string true "The Mechanism name"
// @Success 	204
// @Router 		/mechanism/{name} [delete]
func (controller *MechanismController) Delete(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	// Delete by id
	if err := controller.service.Delete(name); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusNoContent, name)
}
