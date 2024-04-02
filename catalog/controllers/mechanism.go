package controllers

import (
	"net/http"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/model"
	"github.com/FranciscoBarao/catalog/services"
	"github.com/FranciscoBarao/catalog/utils"
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
// @Tags 	mechanisms
// @Produce 	json
// @Param 		data body model.Mechanism true "The Mechanism name"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} model.Mechanism
// @Router 		/mechanism [post]
func (controller *MechanismController) Create(w http.ResponseWriter, r *http.Request) {
	// Deserialize Mechanism input
	var mechanism = &model.Mechanism{}
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

	if err := render.New().JSON(w, http.StatusOK, mechanism); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// Get Mechanisms godoc
// @Summary 	Fetches all Mechanisms
// @Tags 	mechanisms
// @Produce 	json
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
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
// @Success 	200 {object} model.Mechanism
// @Router 		/mechanism/{name} [get]
func (controller *MechanismController) Get(w http.ResponseWriter, r *http.Request) {
	name := utils.GetFieldFromURL(r, "name")

	mechanism, err := controller.service.Get(name)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, mechanism); err != nil {
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
	if err := controller.service.Delete(name); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusNoContent, name); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}
