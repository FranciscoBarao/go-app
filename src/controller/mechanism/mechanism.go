package mechanism

import (
	"net/http"

	"go-app/model"
	"go-app/repository/mechanismRepo"
	"go-app/utils"
)

type repository interface {
	Create(mechanism model.Mechanism) error
	GetAll(sort string) ([]model.Mechanism, error)
	Get(name string) (model.Mechanism, error)
	Delete(mechanism model.Mechanism) error
}

type Controller struct {
	repo repository
}

// InitController initializes the mechanism controller.
func InitController(mechanismRepo *mechanismRepo.MechanismRepository) *Controller {
	return &Controller{
		repo: mechanismRepo,
	}
}

// Create Mechanism godoc
// @Summary 	Creates a Mechanism using a name
// @Mechanisms 	mechanisms
// @Produce 	json
// @Param 		data body model.Mechanism true "The Mechanism name"
// @Success 	200 {object} model.Mechanism
// @Router 		/mechanism/{name} [post]
func (controller *Controller) Create(w http.ResponseWriter, r *http.Request) {

	var mechanism model.Mechanism
	err := utils.DecodeJSONBody(w, r, &mechanism)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	err = controller.repo.Create(mechanism)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &mechanism, http.StatusOK, nil)
}

// Get Mechanisms godoc
// @Summary 	Fetches all Mechanisms
// @Mechanisms 	mechanisms
// @Produce 	json
// @Success 	200 {object} model.Mechanism
// @Router 		/mechanism [get]
func (controller *Controller) GetAll(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sortBy")
	sort, err := utils.GetSort(model.Mechanism{}, sortBy)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	mechanisms, err := controller.repo.GetAll(sort)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &mechanisms, http.StatusOK, nil)
}

// Get Mechanism godoc
// @Summary 	Fetches a specific Mechanism using a name
// @Mechanisms 	mechanisms
// @Produce 	json
// @Param 		name path string true "The Mechanism name"
// @Success 	200 {object} model.Mechanism
// @Router 		/mechanism/{name} [get]
func (controller *Controller) Get(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	mechanism, err := controller.repo.Get(name)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &mechanism, http.StatusOK, nil)
}

// Delete Mechanism godoc
// @Summary 	Deletes a specific Mechanism
// @Mechanisms 	mechanisms
// @Produce 	json
// @Param 		name path string true "The Mechanism name"
// @Success 	200
// @Router 		/mechanism/{name} [delete]
func (controller *Controller) Delete(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	// Get Mechanism by name
	mechanism, err := controller.repo.Get(name)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	// Delete by id
	err = controller.repo.Delete(mechanism)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}
	utils.HTTPHandler(w, name, http.StatusNoContent, nil)
}
