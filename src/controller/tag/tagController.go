package tag

import (
	"errors"
	"net/http"

	"github.com/unrolled/render"

	"go-app/model"
	"go-app/repository/tagRepo"
	"go-app/utils"
)

type repository interface {
	Create(tag model.Tag) error
	GetAll() ([]model.Tag, error)
	Get(name string) (model.Tag, error)
	Delete(tag model.Tag) error
}

type Controller struct {
	repo repository
}

// InitController initializes the boargame controller.
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
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = controller.repo.Create(tag)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
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
func (controller *Controller) GetAll(w http.ResponseWriter, r *http.Request) {

	tags, err := controller.repo.GetAll()
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
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
func (controller *Controller) Get(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	tag, err := controller.repo.Get(name)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	render.New().JSON(w, http.StatusOK, tag)
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
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Delete by id
	err = controller.repo.Delete(tag)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	render.New().JSON(w, http.StatusNoContent, name)
}
