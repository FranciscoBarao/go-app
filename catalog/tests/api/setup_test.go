package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/route"
	"github.com/FranciscoBarao/catalog/internal/transport"
)

//go:generate mockgen -package tests -destination database_mock.go -source setup_test.go

// Database defines the persistence operations needed by the integration tests.
type Database interface {
	CreateCategory(ctx context.Context, c *category.Category) error
	GetCategoryBySlug(ctx context.Context, slug string) (category.Category, error)
	GetCategoryIDBySlug(ctx context.Context, slug string) (uint, error)
	GetAllCategories(ctx context.Context, query listopt.Query) ([]category.Category, int, error)
	DeleteCategory(ctx context.Context, slug string, hard bool) error

	CreateMechanism(ctx context.Context, m *mechanism.Mechanism) error
	GetMechanismBySlug(ctx context.Context, slug string) (mechanism.Mechanism, error)
	GetMechanismIDBySlug(ctx context.Context, slug string) (uint, error)
	GetAllMechanisms(ctx context.Context, query listopt.Query) ([]mechanism.Mechanism, int, error)
	DeleteMechanism(ctx context.Context, slug string, hard bool) error

	CreateContributor(ctx context.Context, c *contributor.Contributor) error
	UpdateContributor(ctx context.Context, c *contributor.Contributor) error
	GetContributorBySlug(ctx context.Context, slug string) (contributor.Contributor, error)
	GetContributorIDBySlug(ctx context.Context, slug string) (uint, error)
	GetAllContributors(ctx context.Context, query listopt.Query) ([]contributor.Contributor, int, error)
	DeleteContributor(ctx context.Context, slug string, hard bool) error

	CreateBoardgame(ctx context.Context, input boardgame.CreateBoardgameDTO) (boardgame.Boardgame, error)
	GetBoardgameByID(ctx context.Context, id uint) (boardgame.Boardgame, error)
	GetBoardgameBySlug(ctx context.Context, slug string) (boardgame.Boardgame, error)
	GetAllBoardgames(ctx context.Context, query listopt.Query, includeDeleted bool) ([]boardgame.Boardgame, int, error)
	UpdateBoardgame(ctx context.Context, bg *boardgame.Boardgame) error
	UpdateBoardgameWithAssociations(ctx context.Context, bg *boardgame.Boardgame, assoc boardgame.UpdateAssociations) error
	DeleteBoardgame(ctx context.Context, id uint, hard bool) error
}

const oauthKey = "secret-key"

// assertEnvelope returns an apitest assertion that decodes the paginated
// response envelope and checks the page window and item/data counts.
func assertEnvelope(t *testing.T, wantLen, wantPage, wantPageSize, wantTotal, wantTotalPages int) func(*http.Response, *http.Request) error {
	return func(res *http.Response, _ *http.Request) error {
		var body struct {
			Data       []json.RawMessage `json:"data"`
			Page       int               `json:"page"`
			PageSize   int               `json:"pageSize"`
			TotalItems int               `json:"totalItems"`
			TotalPages int               `json:"totalPages"`
		}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return err
		}
		require.Len(t, body.Data, wantLen)
		require.Equal(t, wantPage, body.Page)
		require.Equal(t, wantPageSize, body.PageSize)
		require.Equal(t, wantTotal, body.TotalItems)
		require.Equal(t, wantTotalPages, body.TotalPages)
		return nil
	}
}

type Base struct {
	router      *chi.Mux
	oauthHeader string
	dbMock      *MockDatabase
}

// NewBase prepares the test environment.
func NewBase(t *testing.T) *Base {
	log := logging.FromCtx(context.Background())
	log.Debug().Msg("setup starting..")

	// Mirrors main.go: QUERY is not in chi's built-in method table, so it has to
	// be registered before the routes below are attached.
	// TODO: drop once QUERY leaves IETF draft and chi supports it natively.
	chi.RegisterMethod("QUERY")

	ctrl := gomock.NewController(t)
	mock := NewMockDatabase(ctrl)

	categorySvc := category.NewService(mock)
	mechanismSvc := mechanism.NewService(mock)
	contributorSvc := contributor.NewService(mock)
	boardgameSvc := boardgame.NewService(mock, categorySvc, mechanismSvc, contributorSvc)

	bgController := transport.NewBoardgameController(boardgameSvc)
	categoryController := transport.NewCategoryController(categorySvc)
	mechanismController := transport.NewMechanismController(mechanismSvc)
	contributorController := transport.NewContributorController(contributorSvc)

	router := chi.NewRouter()
	route.AddBoardGameRouter(router, oauthKey, bgController)
	route.AddCategoryRouter(router, oauthKey, categoryController)
	route.AddMechanismRouter(router, oauthKey, mechanismController)
	route.AddContributorRouter(router, oauthKey, contributorController)

	log.Debug().Msg("setup complete")
	return &Base{
		router:      router,
		oauthHeader: oauthKey,
		dbMock:      mock,
	}
}
