package tests

//go:generate mockgen -package tests -destination mock_database_test.go -source setup_test.go

import (
	"context"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/boardgame"
	"github.com/FranciscoBarao/catalog/category"
	"github.com/FranciscoBarao/catalog/mechanism"
	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/route"
	"github.com/FranciscoBarao/catalog/tag"
	"github.com/FranciscoBarao/catalog/transport"
)

// Database defines the persistence operations needed by the integration tests.
type Database interface {
	Create(value interface{}) error
	Read(value interface{}, sort, search, identifier string) error
	Update(value interface{}) error
	Delete(value interface{}) error
	ReplaceAssociatons(model interface{}, association string, values interface{}) error
}

const oauthKey = "secret-key"

type Base struct {
	router      *chi.Mux
	oauthHeader string
	dbMock      *MockDatabase
}

// Prepares test environment
func NewBase(t *testing.T) *Base {
	log := middleware.FromCtx(context.Background())
	log.Debug().Msg("setup starting..")

	// Setup database mock
	ctrl := gomock.NewController(t)
	mock := NewMockDatabase(ctrl)

	// Initialize Services
	tagSvc := tag.NewService(mock)
	categorySvc := category.NewService(mock)
	mechanismSvc := mechanism.NewService(mock)
	boardgameSvc := boardgame.NewService(mock, tagSvc, categorySvc, mechanismSvc)

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
