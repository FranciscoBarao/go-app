package tag

import (
	"net/http"

	"go-app/model"
	"go-app/repository/tagRepo"
	"go-app/utils"
)

type repository interface {
	Create(tag model.Tag) error
	GetAll(sort string) ([]model.Tag, error)
	Get(name string) (model.Tag, error)
	Delete(tag model.Tag) error
}

type Controller struct {
	repo repository
}

// InitController initializes the tag controller.
func InitController(tagRepo *tagRepo.TagRepository) *Controller {
	return &Controller{
		repo: tagRepo,
	}
}

// Create Tag godoc
// @Summary 	Creates a Tag using a name
// @Tags 		tags
// @Produce 	json
// @Param 		data body model.Tag true "The Tag name"
// @Success 	200 {object} model.Tag
// @Router 		/tag/{name} [post]
func (controller *Controller) Create(w http.ResponseWriter, r *http.Request) {

	var tag model.Tag
	err := utils.DecodeJSONBody(w, r, &tag)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	err = controller.repo.Create(tag)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &tag, http.StatusOK, nil)
}

// Get Tags godoc
// @Summary 	Fetches all Tags
// @Tags 		tags
// @Produce 	json
// @Success 	200 {object} model.Tag
// @Router 		/tag [get]
func (controller *Controller) GetAll(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sortBy")
	sort, err := utils.GetSort(model.Tag{}, sortBy)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	tags, err := controller.repo.GetAll(sort)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &tags, http.StatusOK, nil)
}

// Get Tag godoc
// @Summary 	Fetches a specific Tag using a name
// @Tags 		tags
// @Produce 	json
// @Param 		name path string true "The Tag name"
// @Success 	200 {object} model.Tag
// @Router 		/tag/{name} [get]
func (controller *Controller) Get(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	tag, err := controller.repo.Get(name)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	utils.HTTPHandler(w, &tag, http.StatusOK, nil)
}

// Delete Tag godoc
// @Summary 	Deletes a specific Tag
// @Tags 		tags
// @Produce 	json
// @Param 		name path string true "The Tag name"
// @Success 	200
// @Router 		/tag/{name} [delete]
func (controller *Controller) Delete(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	// Get Tag by name
	tag, err := controller.repo.Get(name)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}

	// Delete by id
	err = controller.repo.Delete(tag)
	if err != nil {
		utils.HTTPHandler(w, nil, 0, err)
		return
	}
	utils.HTTPHandler(w, name, http.StatusNoContent, nil)
}
