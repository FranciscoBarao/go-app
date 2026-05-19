package tests

import (
	"context"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/route"
	"github.com/FranciscoBarao/catalog/internal/tag"
	"github.com/FranciscoBarao/catalog/internal/transport"
)

//go:generate mockgen -package tests -destination database_mock.go -source setup_test.go

// Database defines the persistence operations needed by the integration tests.
// It combines all four resource Database interfaces into a single mock target.
type Database interface {
	// Tag methods
	CreateTag(ctx context.Context, t *tag.Tag) error
	GetTag(ctx context.Context, name string) (tag.Tag, error)
	GetAllTags(ctx context.Context, filter listopt.Params) ([]tag.Tag, error)
	DeleteTag(ctx context.Context, name string) error

	// Category methods
	CreateCategory(ctx context.Context, c *category.Category) error
	GetCategory(ctx context.Context, name string) (category.Category, error)
	GetAllCategories(ctx context.Context, filter listopt.Params) ([]category.Category, error)
	DeleteCategory(ctx context.Context, name string) error

	// Mechanism methods
	CreateMechanism(ctx context.Context, m *mechanism.Mechanism) error
	GetMechanism(ctx context.Context, name string) (mechanism.Mechanism, error)
	GetAllMechanisms(ctx context.Context, filter listopt.Params) ([]mechanism.Mechanism, error)
	DeleteMechanism(ctx context.Context, name string) error

	// Boardgame methods
	CreateBoardgame(ctx context.Context, bg *boardgame.Boardgame) error
	GetBoardgameByID(ctx context.Context, id uint) (boardgame.Boardgame, error)
	GetAllBoardgames(ctx context.Context, filter listopt.Params) ([]boardgame.Boardgame, error)
	UpdateBoardgame(ctx context.Context, bg *boardgame.Boardgame) error
	DeleteBoardgame(ctx context.Context, id uint) error
	ReplaceBoardgameTags(ctx context.Context, boardgameID uint, tags []tag.Tag) error
}

const oauthKey = "secret-key"

type Base struct {
	router      *chi.Mux
	oauthHeader string
	dbMock      *MockDatabase
}

// Prepares test environment
func NewBase(t *testing.T) *Base {
	log := logging.FromCtx(context.Background())
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
