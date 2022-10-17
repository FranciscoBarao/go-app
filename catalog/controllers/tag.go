package controllers

import (
	"net/http"

	"catalog/middleware"
	"catalog/model"
	"catalog/repositories"
	"catalog/utils"

	"github.com/unrolled/render"
)

type tagRepository interface {
	Create(tag model.Tag) error
	GetAll(sort string) ([]model.Tag, error)
	Get(name string) (model.Tag, error)
	Delete(tag model.Tag) error
}

type TagController struct {
	repo tagRepository
}

// InitController initializes the tag controller.
func InitTagController(tagRepo *repositories.TagRepository) *TagController {
	return &TagController{
		repo: tagRepo,
	}
}

// Create Tag godoc
// @Summary 	Creates a Tag using a name
// @Tags 		tags
// @Produce 	json
// @Param 		data body model.Tag true "The Tag name"
// @Success 	200 {object} model.Tag
// @Router 		/tag [post]
func (controller *TagController) Create(w http.ResponseWriter, r *http.Request) {

	// Deserialize Tag input
	var tag model.Tag
	if err := utils.DecodeJSONBody(w, r, &tag); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate Tag input
	if err := utils.ValidateStruct(&tag); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.repo.Create(tag); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, tag)
}

// Get Tags godoc
// @Summary 	Fetches all Tags
// @Tags 		tags
// @Produce 	json
// @Success 	200 {object} model.Tag
// @Router 		/tag [get]
func (controller *TagController) GetAll(w http.ResponseWriter, r *http.Request) {

	sortBy := r.URL.Query().Get("sortBy")
	sort, err := utils.GetSort(model.Tag{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	tags, err := controller.repo.GetAll(sort)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	render.New().JSON(w, http.StatusOK, tags)
}

// Get Tag godoc
// @Summary 	Fetches a specific Tag using a name
// @Tags 		tags
// @Produce 	json
// @Param 		name path string true "The Tag name"
// @Success 	200 {object} model.Tag
// @Router 		/tag/{name} [get]
func (controller *TagController) Get(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")
	tag, err := controller.repo.Get(name)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusOK, tag)
}

// Delete Tag godoc
// @Summary 	Deletes a specific Tag
// @Tags 		tags
// @Produce 	json
// @Param 		name path string true "The Tag name"
// @Success 	204
// @Router 		/tag/{name} [delete]
func (controller *TagController) Delete(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	// Get Tag by name
	tag, err := controller.repo.Get(name)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Delete by id
	if err := controller.repo.Delete(tag); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	render.New().JSON(w, http.StatusNoContent, name)
}
