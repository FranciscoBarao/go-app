package controllers

import (
	"net/http"

	"catalog/middleware"
	"catalog/model"
	"catalog/repositories"
	"catalog/utils"

	"github.com/unrolled/render"
)

type mechanismRepository interface {
	Create(mechanism model.Mechanism) error
	GetAll(sort string) ([]model.Mechanism, error)
	Get(name string) (model.Mechanism, error)
	Delete(mechanism model.Mechanism) error
}

type MechanismController struct {
	repo mechanismRepository
}

// InitController initializes the mechanism controller.
func InitMechanismController(mechanismRepo *repositories.MechanismRepository) *MechanismController {
	return &MechanismController{
		repo: mechanismRepo,
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
	var mechanism model.Mechanism
	if err := utils.DecodeJSONBody(w, r, &mechanism); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate Mechanism input
	if err := utils.ValidateStruct(&mechanism); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.repo.Create(mechanism); err != nil {
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

	mechanisms, err := controller.repo.GetAll(sort)
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

	mechanism, err := controller.repo.Get(name)
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

	// Get Mechanism by name
	mechanism, err := controller.repo.Get(name)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Delete by id
	if err := controller.repo.Delete(mechanism); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusNoContent, name)
}
