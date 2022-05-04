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
	GetByName(name string) (model.Tag, error)
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

// Method that Creates a tag based on json input
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

// Method that Gets a tag based on a unique name
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

// Method that Gets a tag based on a unique name
func (controller *Controller) GetByName(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	tag, err := controller.repo.GetByName(name)
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

// Method that Deletes a tag based on the unique name
func (controller *Controller) Delete(w http.ResponseWriter, r *http.Request) {

	name := utils.GetFieldFromURL(r, "name")

	// Get Tag by name
	tag, err := controller.repo.GetByName(name)
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
