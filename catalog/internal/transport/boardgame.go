package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/utils"
)

//go:generate mockgen -package transport -destination boardgame_mock.go . BoardgameService

// BoardgameService defines the interface for boardgame business logic.
type BoardgameService interface {
	Create(ctx context.Context, req *boardgame.CreateBoardgameRequest, parentSlug string) (boardgame.Boardgame, error)
	GetAll(ctx context.Context, includeDeleted bool, query listopt.Query) ([]boardgame.Boardgame, int, error)
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
// @Summary 	Creates a Boardgame (or an expansion when a parent slug is given)
// @Tags 		boardgames
// @Produce 	json
// @Param 		slug path string false "Parent boardgame slug (when creating an expansion)"
// @Param 		data body boardgame.CreateBoardgameRequest true "Boardgame"
// @Success 	200 {object} boardgame.Boardgame
// @Router 		/boardgame [post]
// @Router 		/boardgame/{slug}/expansion [post]
func (controller *BoardgameController) Create(w http.ResponseWriter, r *http.Request) {
	var req boardgame.CreateBoardgameRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		switch {
		case errors.Is(err, utils.ErrBadRequest):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, utils.ErrRequestTooLarge):
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("decode create boardgame")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := utils.ValidateStruct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parentSlug := utils.GetFieldFromURL(r, "parentSlug")
	if parentSlug == "" {
		parentSlug = utils.GetFieldFromURL(r, "slug")
	}

	bg, err := controller.service.Create(r.Context(), &req, parentSlug)
	if err != nil {
		switch {
		case errors.Is(err, boardgame.ErrValidation):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, boardgame.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, boardgame.ErrAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, boardgame.ErrMissingReference):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, boardgame.ErrParentIsExpansion):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("create boardgame")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := render.New().JSON(w, http.StatusOK, bg); err != nil {
		logging.FromCtx(r.Context()).Error().Err(err).Msg("encode boardgame")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetAll Boardgames godoc
// @Summary 	Lists Boardgames with pagination, sorting, and filtering
// @Tags 		boardgames
// @Produce 	json
// @Param 		page query int false "Page number (default 1)"
// @Param 		pageSize query int false "Items per page (default 10, max 100)"
// @Param 		sort query string false "Sort as field.order, e.g. name.asc"
// @Param 		filter query []string false "Filter as field.op.value, e.g. min_players.ge.3; repeatable, AND-combined"
// @Param 		include_deleted query bool false "Include soft-deleted boardgames"
// @Success 	200 {object} BoardgamePage
// @Router 		/boardgame [get]
func (controller *BoardgameController) GetAll(w http.ResponseWriter, r *http.Request) {
	req, err := newQueryRequestFromURL(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	controller.list(w, r, req)
}

// Query lists Boardgames with filtering, sorting, and pagination.
//
// It handles HTTP QUERY /api/boardgame with a JSON body (pagination, sort,
// filters, include_deleted), for filters too complex to express as GET query
// parameters. Not represented in Swagger because OpenAPI 2.0 has no QUERY
// method; see GetAll for the same envelope.
func (controller *BoardgameController) Query(w http.ResponseWriter, r *http.Request) {
	var req QueryRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		switch {
		case errors.Is(err, utils.ErrBadRequest):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, utils.ErrRequestTooLarge):
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("decode boardgame query")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	controller.list(w, r, req)
}

// list resolves a validated list request and writes the paginated envelope,
// shared by the GET and QUERY entry points.
func (controller *BoardgameController) list(w http.ResponseWriter, r *http.Request, req QueryRequest) {
	query, err := req.toQuery(boardgame.QueryAllowlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	boardgames, total, err := controller.service.GetAll(r.Context(), req.IncludeDeleted, query)
	if err != nil {
		logging.FromCtx(r.Context()).Error().Err(err).Msg("list boardgames")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, newPaginatedResponse(boardgames, total, query.Pagination)); err != nil {
		logging.FromCtx(r.Context()).Error().Err(err).Msg("encode boardgames")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, fmt.Sprintf("invalid boardgame id %s", idStr), http.StatusBadRequest)
		return
	}

	bg, err := controller.service.GetByID(r.Context(), uint(parsedID))
	if err != nil {
		switch {
		case errors.Is(err, boardgame.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("get boardgame by id")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := render.New().JSON(w, http.StatusOK, bg); err != nil {
		logging.FromCtx(r.Context()).Error().Err(err).Msg("encode boardgame")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		switch {
		case errors.Is(err, boardgame.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("get boardgame")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := render.New().JSON(w, http.StatusOK, bg); err != nil {
		logging.FromCtx(r.Context()).Error().Err(err).Msg("encode boardgame")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	if err := decodeJSONBody(w, r, &req); err != nil {
		switch {
		case errors.Is(err, utils.ErrBadRequest):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, utils.ErrRequestTooLarge):
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("decode update boardgame")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := utils.ValidateStruct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slugStr := utils.GetFieldFromURL(r, "slug")
	err := controller.service.Update(r.Context(), &req, slugStr)
	if err != nil {
		switch {
		case errors.Is(err, boardgame.ErrValidation):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, boardgame.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, boardgame.ErrAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, boardgame.ErrMissingReference):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("update boardgame")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := render.New().JSON(w, http.StatusOK, req); err != nil {
		logging.FromCtx(r.Context()).Error().Err(err).Msg("encode boardgame")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	err := controller.service.DeleteBySlug(r.Context(), slugStr, hard)
	if err != nil {
		switch {
		case errors.Is(err, boardgame.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("delete boardgame")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
	if err := decodeJSONBody(w, r, rating); err != nil {
		switch {
		case errors.Is(err, utils.ErrBadRequest):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, utils.ErrRequestTooLarge):
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("decode boardgame rating")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := utils.ValidateStruct(rating); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slugStr := utils.GetFieldFromURL(r, "slug")
	user, err := utils.GetUsernameFromToken(r)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorized):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("boardgame rating user")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := controller.service.Rate(r.Context(), rating, slugStr, user); err != nil {
		switch {
		case errors.Is(err, boardgame.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			logging.FromCtx(r.Context()).Error().Err(err).Msg("rate boardgame")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := render.New().JSON(w, http.StatusOK, rating); err != nil {
		logging.FromCtx(r.Context()).Error().Err(err).Msg("encode boardgame rating")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
