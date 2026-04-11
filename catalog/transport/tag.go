package transport

import (
	"context"
	"net/http"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/tag"
	"github.com/FranciscoBarao/catalog/utils"
)

//go:generate mockgen -package transport -destination tag_mock.go . TagService

// TagService defines the interface for tag business logic.
type TagService interface {
	Create(ctx context.Context, t *tag.Tag) error
	GetAll(ctx context.Context, sort string) ([]tag.Tag, error)
	Get(ctx context.Context, name string) (tag.Tag, error)
	Delete(ctx context.Context, name string) error
}

// TagController handles HTTP requests for tag operations.
type TagController struct {
	service TagService
}

// NewTagController initializes the tag controller.
func NewTagController(tagSvc TagService) *TagController {
	return &TagController{
		service: tagSvc,
	}
}

// Create Tag godoc
// @Summary 	Creates a Tag using a name
// @Tags 		tags
// @Produce 	json
// @Param 		data body tag.Tag true "The Tag name"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} tag.Tag
// @Router 		/tag [post]
func (controller *TagController) Create(w http.ResponseWriter, r *http.Request) {
	// Deserialize Tag input
	var t = &tag.Tag{}
	if err := utils.DecodeJSONBody(w, r, t); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate Tag input
	if err := utils.ValidateStruct(t); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.service.Create(context.Background(), t); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, t); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// GetAll Tags godoc
// @Summary 	Fetches all Tags
// @Tags 		tags
// @Produce 	json
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} tag.Tag
// @Router 		/tag [get]
func (controller *TagController) GetAll(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	sort, err := utils.GetSort(tag.Tag{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	tags, err := controller.service.GetAll(context.Background(), sort)
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
// @Success 	200 {object} tag.Tag
// @Router 		/tag/{name} [get]
func (controller *TagController) Get(w http.ResponseWriter, r *http.Request) {
	name := utils.GetFieldFromURL(r, "name")
	tag, err := controller.service.Get(context.Background(), name)
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
	if err := controller.service.Delete(context.Background(), name); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusNoContent, name); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}
