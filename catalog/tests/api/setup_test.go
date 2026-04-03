package tests

import (
	"context"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/boardgame"
	"github.com/FranciscoBarao/catalog/category"
	"github.com/FranciscoBarao/catalog/database"
	"github.com/FranciscoBarao/catalog/mechanism"
	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/route"
	"github.com/FranciscoBarao/catalog/tag"
	"github.com/FranciscoBarao/catalog/transport"
)

const oauthKey = "secret-key"

type Base struct {
	router      *chi.Mux
	oauthHeader string
	dbMock      *database.MockDatabase
}

// Prepares test environment
func NewBase(t *testing.T) *Base {
	log := middleware.FromCtx(context.Background())
	log.Debug().Msg("setup starting..")

	// Setup database mock
	ctrl := gomock.NewController(t)
	mock := database.NewMockDatabase(ctrl)

	// Initialize Services
	tagSvc := tag.NewTagService(mock)
	categorySvc := category.NewCategoryService(mock)
	mechanismSvc := mechanism.NewMechanismService(mock)
	boardgameSvc := boardgame.NewBoardgameService(mock, tagSvc, categorySvc, mechanismSvc)

	// Initialize Controllers
	bgController := transport.NewBoardgameController(boardgameSvc)
	tagController := transport.NewTagController(tagSvc)
	categoryController := transport.NewCategoryController(categorySvc)
	mechanismController := transport.NewMechanismController(mechanismSvc)

	// Adds Routers
	router := chi.NewRouter()
	route.AddBoardGameRouter(router, oauthKey, bgController)
	route.AddTagRouter(router, oauthKey, tagController)
	route.AddCategoryRouter(router, oauthKey, categoryController)
	route.AddMechanismRouter(router, oauthKey, mechanismController)

	log.Debug().Msg("setup complete")
	return &Base{
		router:      router,
		oauthHeader: oauthKey,
		dbMock:      mock,
	}
}
