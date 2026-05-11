package transport

import (
	"context"
	"net/http"
	"strconv"

	"github.com/unrolled/render"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/query"
	"github.com/FranciscoBarao/catalog/internal/utils"
)

//go:generate mockgen -package transport -destination boardgame_mock.go . BoardgameService

// BoardgameService defines the interface for boardgame business logic.
type BoardgameService interface {
	Create(ctx context.Context, bg *boardgame.Boardgame, id uint) error
	GetAll(ctx context.Context, opts ...query.Option) ([]boardgame.Boardgame, error)
	GetByID(ctx context.Context, id uint) (boardgame.Boardgame, error)
	Update(ctx context.Context, req *boardgame.UpdateBoardgameRequest, id uint) error
	DeleteByID(ctx context.Context, id uint) error
	Rate(ctx context.Context, rating *boardgame.Rating, id uint, username string) error
}

// BoardgameController handles HTTP requests for boardgame operations.
type BoardgameController struct {
	service BoardgameService
}

// NewBoardgameController initializes the boardgame and the associations controller
func NewBoardgameController(boardGameSvc BoardgameService) *BoardgameController {
	return &BoardgameController{
		service: boardGameSvc,
	}
}

// Create Boardgame godoc
// @Summary 	Creates a Boardgame based on a json body
// @Tags 		boardgames
// @Produce 	json
// @Param 		data body boardgame.Boardgame true "The input Boardgame struct"
// @Param 		id path int false "The Boardgame id indicating this is an Expansion"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} boardgame.Boardgame
// @Router 		/boardgame [post]
func (controller *BoardgameController) Create(w http.ResponseWriter, r *http.Request) {
	// Deserialize into request
	var req boardgame.CreateBoardgameRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate input
	if err := utils.ValidateStruct(&req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Construct Boardgame
	bg := boardgame.NewBoardgame(&req)

	// Get Id from url - If its an expansion
	var parentID uint
	if id := utils.GetFieldFromURL(r, "id"); id != "" {
		parsedID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			middleware.ErrorHandler(w, middleware.NewError(http.StatusBadRequest, "Invalid boardgame ID: "+id))
			return
		}
		parentID = uint(parsedID)
	}

	if err := controller.service.Create(context.Background(), bg, parentID); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, bg); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// GetAll Boardgames godoc
// @Summary 	Fetches all Boardgames
// @Tags 		boardgames
// @Produce 	json
// @Success 	200 {object} boardgame.Boardgame
// @Router 		/boardgame [get]
func (controller *BoardgameController) GetAll(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	col, order, err := utils.GetSort(boardgame.Boardgame{}, sortBy)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	var opts []query.Option
	if col != "" {
		opts = append(opts, query.WithSort(col, order))
	}

	/*
		TODO - Reevaluate and rething this filter strategy

		// @Param  filterBy query string  false  "Filter using field.value (For String partial find) OR field.operator.value"

		filterBy := r.URL.Query().Get("filterBy")
		filterBody, filterValue, err := utils.GetFilters(boardgame.Boardgame{}, filterBy)
		if err != nil {
			middleware.ErrorHandler(w, err)
			return
		}
	*/

	boardgames, err := controller.service.GetAll(context.Background(), opts...)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, boardgames); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// Get Boardgame by id godoc
// @Summary 	Fetches a specific Boardgame using an id
// @Tags 		boardgames
// @Produce 	json
// @Param 		id path int true "The Boardgame unique id"
// @Success 	200 {object} boardgame.Boardgame
// @Router 		/boardgame/{id} [get]
func (controller *BoardgameController) Get(w http.ResponseWriter, r *http.Request) {
	id := utils.GetFieldFromURL(r, "id")
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		err = middleware.NewError(http.StatusBadRequest, "Invalid boardgame ID: "+id)
		middleware.ErrorHandler(w, err)
		return
	}

	boardgame, err := controller.service.GetByID(context.Background(), uint(parsedID))
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
	if err := render.New().JSON(w, http.StatusOK, boardgame); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// Update Boardgame by id godoc
// @Summary 	Updates a specific Boardgame via Id
// @Tags 		boardgames
// @Produce 	json
// @Param 		id path int true "The Boardgame id"
// @Param 		data body UpdateBoardgameRequest true "The partial Boardgame fields to update"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} boardgame.Boardgame
// @Router 		/boardgame/{id} [patch]
func (controller *BoardgameController) Update(w http.ResponseWriter, r *http.Request) {
	// Deserialize into request
	var req boardgame.UpdateBoardgameRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate input
	if err := utils.ValidateStruct(&req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Get ID from URL
	id := utils.GetFieldFromURL(r, "id")
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		middleware.ErrorHandler(w, middleware.NewError(http.StatusBadRequest, "Invalid boardgame ID: "+id))
		return
	}

	// Updates Boardgame
	if err := controller.service.Update(context.Background(), &req, uint(parsedID)); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, req); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// Delete Boardgame by id godoc
// @Summary 	Deletes a specific Boardgame via Id
// @Tags 		boardgames
// @Produce 	json
// @Param 		id path int true "The Boardgame id"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	204
// @Router 		/boardgame/{id} [delete]
func (controller *BoardgameController) Delete(w http.ResponseWriter, r *http.Request) {
	id := utils.GetFieldFromURL(r, "id")
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		err = middleware.NewError(http.StatusBadRequest, "Invalid boardgame ID: "+id)
		middleware.ErrorHandler(w, err)
		return
	}

	// Delete by Id
	if err := controller.service.DeleteByID(context.Background(), uint(parsedID)); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusNoContent, id); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}

// Rate a Boardgame godoc
// @Summary 	Rates a Boardgame
// @Tags 		boardgames
// @Produce 	json
// @Param 		id path int true "The Boardgame id"
// @Param 		Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 	200 {object} boardgame.Rating
// @Router 		/boardgame/{id}/rate [post]
func (controller *BoardgameController) Rate(w http.ResponseWriter, r *http.Request) {
	// Deserialize Rating input
	var rating = &boardgame.Rating{}
	if err := utils.DecodeJSONBody(w, r, rating); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Validate Boardgame input
	if err := utils.ValidateStruct(rating); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	// Get Boardgame Id from url
	id := utils.GetFieldFromURL(r, "id")
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		err = middleware.NewError(http.StatusBadRequest, "Invalid boardgame ID: "+id)
		middleware.ErrorHandler(w, err)
		return
	}

	// Get username from oauth Token
	user, err := utils.GetUsernameFromToken(r)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := controller.service.Rate(context.Background(), rating, uint(parsedID), user); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	if err := render.New().JSON(w, http.StatusOK, rating); err != nil {
		middleware.ErrorHandler(w, err)
		return
	}
}
