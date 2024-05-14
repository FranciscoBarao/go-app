package controllers

import (
	"net/http"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/model"
	"github.com/FranciscoBarao/catalog/services"
	"github.com/FranciscoBarao/catalog/utils"
)

type tagService interface {
	Create(tag *model.Tag) error
	GetAll(sort string) ([]model.Tag, error)
	Get(name string) (model.Tag, error)
	Delete(name string) error
}

type TagController struct {
	service tagService
}

// InitController initializes the tag controller.
func InitTagController(tagSvc *services.TagService) *TagController {
	return &TagController{
		service: tagSvc,
	}
}

// Create Tag godoc
// @Summary 	Creates a Tag using a name
// @Tags 		tags
// @Produce 	json
// @Param 		data body model.Tag true "The Tag name"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} model.Tag
// @Router 		/tag [post]
func (controller *TagController) Create(w http.ResponseWriter, r *http.Request) {
	// Deserialize Tag input
	var tag = &model.Tag{}
	if err := utils.DecodeJSONBody(w, r, tag); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate Tag input
	if err := utils.ValidateStruct(tag); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.service.Create(tag); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, tag); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// Get Tags godoc
// @Summary 	Fetches all Tags
// @Tags 		tags
// @Produce 	json
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} model.Tag
// @Router 		/tag [get]
func (controller *TagController) GetAll(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	sort, err := utils.GetSort(model.Tag{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	tags, err := controller.service.GetAll(sort)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, tags); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// Get Tag godoc
// @Summary 	Fetches a specific Tag using a name
// @Tags 		tags
// @Produce 	json
// @Param 		name path string true "The Tag name"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} model.Tag
// @Router 		/tag/{name} [get]
func (controller *TagController) Get(w http.ResponseWriter, r *http.Request) {
	name := utils.GetFieldFromURL(r, "name")
	tag, err := controller.service.Get(name)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, tag); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// Delete Tag godoc
// @Summary 	Deletes a specific Tag
// @Tags 		tags
// @Produce 	json
// @Param 		name path string true "The Tag name"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	204
// @Router 		/tag/{name} [delete]
func (controller *TagController) Delete(w http.ResponseWriter, r *http.Request) {
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
