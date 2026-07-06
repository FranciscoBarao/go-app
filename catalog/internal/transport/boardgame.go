package transport

import (
	"context"
	"net/http"
	"strconv"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/utils"
)

//go:generate mockgen -package transport -destination boardgame_mock.go . BoardgameService

// BoardgameService defines the interface for boardgame business logic.
type BoardgameService interface {
	Create(ctx context.Context, req *boardgame.CreateBoardgameRequest, parentSlug string) (boardgame.Boardgame, error)
	GetAll(ctx context.Context, includeDeleted bool, opts ...listopt.Option) ([]boardgame.Boardgame, error)
	GetBySlug(ctx context.Context, slug string) (boardgame.Boardgame, error)
	GetByID(ctx context.Context, id uint) (boardgame.Boardgame, error)
	Update(ctx context.Context, req *boardgame.UpdateBoardgameRequest, slug string) error
	DeleteBySlug(ctx context.Context, slug string, hard bool) error
	Rate(ctx context.Context, rating *boardgame.Rating, slug string, username string) error
}

// BoardgameController handles HTTP requests for boardgame operations.
type BoardgameController struct {
	service BoardgameService
}

// NewBoardgameController initializes the boardgame controller.
func NewBoardgameController(boardGameSvc BoardgameService) *BoardgameController {
	return &BoardgameController{service: boardGameSvc}
}

// Create Boardgame godoc
// @Summary 	Creates a Boardgame
// @Tags 		boardgames
// @Produce 	json
// @Param 		data body boardgame.CreateBoardgameRequest true "Boardgame"
// @Success 	200 {object} boardgame.Boardgame
// @Router 		/boardgame [post]
func (controller *BoardgameController) Create(w http.ResponseWriter, r *http.Request) {
	var req boardgame.CreateBoardgameRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := utils.ValidateStruct(&req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	parentSlug := utils.GetFieldFromURL(r, "parentSlug")
	if parentSlug == "" {
		parentSlug = utils.GetFieldFromURL(r, "slug")
	}

	bg, err := controller.service.Create(r.Context(), &req, parentSlug)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, bg); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// GetAll Boardgames godoc
// @Summary 	Fetches all Boardgames
// @Tags 		boardgames
// @Produce 	json
// @Success 	200 {array} boardgame.Boardgame
// @Router 		/boardgame [get]
func (controller *BoardgameController) GetAll(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	col, order, err := utils.GetSort(boardgame.Boardgame{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	var opts []listopt.Option
	if col != "" {
		opts = append(opts, listopt.WithSort(col, order))
	}

	filterBy := r.URL.Query().Get("filterBy")
	fcol, fop, fval, err := utils.GetFilter(boardgame.Boardgame{}, filterBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if fcol != "" {
		opts = append(opts, listopt.WithFilter(fcol, fop, fval))
	}

	includeDeleted := r.URL.Query().Get("include_deleted") == "true"

	boardgames, err := controller.service.GetAll(r.Context(), includeDeleted, opts...)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, boardgames); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// GetByID Boardgame godoc
// @Summary 	Fetches a Boardgame by numeric id (internal)
// @Tags 		boardgames
// @Produce 	json
// @Param 		id path int true "Boardgame id"
// @Success 	200 {object} boardgame.Boardgame
// @Router 		/boardgame/by-id/{id} [get]
func (controller *BoardgameController) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := utils.GetFieldFromURL(r, "id")
	parsedID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		middleware.ErrorHandler(w, middleware.NewError(http.StatusBadRequest, "Invalid boardgame ID: "+idStr))
		return
	}

	bg, err := controller.service.GetByID(r.Context(), uint(parsedID))
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, bg); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// Get Boardgame by slug godoc
// @Summary 	Fetches a Boardgame by slug
// @Tags 		boardgames
// @Produce 	json
// @Param 		slug path string true "Boardgame slug"
// @Success 	200 {object} boardgame.Boardgame
// @Router 		/boardgame/{slug} [get]
func (controller *BoardgameController) Get(w http.ResponseWriter, r *http.Request) {
	slugStr := utils.GetFieldFromURL(r, "slug")
	bg, err := controller.service.GetBySlug(r.Context(), slugStr)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, bg); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// Update Boardgame godoc
// @Summary 	Updates a Boardgame
// @Tags 		boardgames
// @Produce 	json
// @Param 		slug path string true "Boardgame slug"
// @Param 		data body boardgame.UpdateBoardgameRequest true "Fields to update"
// @Success 	200 {object} boardgame.UpdateBoardgameRequest
// @Router 		/boardgame/{slug} [patch]
func (controller *BoardgameController) Update(w http.ResponseWriter, r *http.Request) {
	var req boardgame.UpdateBoardgameRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := utils.ValidateStruct(&req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	slugStr := utils.GetFieldFromURL(r, "slug")
	if err := controller.service.Update(r.Context(), &req, slugStr); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, req); err != nil {
		middleware.ErrorHandler(w, err)
	}
}

// Delete Boardgame godoc
// @Summary 	Deletes a Boardgame
// @Tags 		boardgames
// @Produce 	json
// @Param 		slug path string true "Boardgame slug"
// @Param 		hard query bool false "Hard delete"
// @Success 	204
// @Router 		/boardgame/{slug} [delete]
func (controller *BoardgameController) Delete(w http.ResponseWriter, r *http.Request) {
	slugStr := utils.GetFieldFromURL(r, "slug")
	hard := r.URL.Query().Get("hard") == "true"

	if err := controller.service.DeleteBySlug(r.Context(), slugStr, hard); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Rate a Boardgame godoc
// @Summary 	Rates a Boardgame
// @Tags 		boardgames
// @Produce 	json
// @Param 		slug path string true "Boardgame slug"
// @Success 	200 {object} boardgame.Rating
// @Router 		/boardgame/{slug}/rate [post]
func (controller *BoardgameController) Rate(w http.ResponseWriter, r *http.Request) {
	var rating = &boardgame.Rating{}
	if err := utils.DecodeJSONBody(w, r, rating); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := utils.ValidateStruct(rating); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	slugStr := utils.GetFieldFromURL(r, "slug")
	user, err := utils.GetUsernameFromToken(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.service.Rate(r.Context(), rating, slugStr, user); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, rating); err != nil {
		middleware.ErrorHandler(w, err)
	}
}
