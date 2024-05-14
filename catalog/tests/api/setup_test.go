package tests

import (
	"context"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"

	"github.com/FranciscoBarao/catalog/controllers"
	"github.com/FranciscoBarao/catalog/middleware/logging"
	"github.com/FranciscoBarao/catalog/repositories"
	"github.com/FranciscoBarao/catalog/route"
	"github.com/FranciscoBarao/catalog/services"
)

const oauthKey = "secret-key"

type Base struct {
	router      *chi.Mux
	oauthHeader string
	dbMock      *repositories.MockDatabase
}

// Prepares test environment
func NewBase(t *testing.T) *Base {
	log := logging.FromCtx(context.Background())
	log.Debug().Msg("setup starting..")

	// Setup database mock
	mock := repositories.NewMockDatabase(gomock.NewController(t))

	// Fetch Oauth Key
	//oauthKey, _ := os.LookupEnv("OAUTH_KEY")

	// Set Repositories & Controllers & Services
	repositories := repositories.InitRepositories(mock)
	services := services.InitServices(repositories)
	controllers := controllers.InitControllers(services)

	// Adds Routers
	router := chi.NewRouter()
	route.AddBoardGameRouter(router, oauthKey, controllers.BoardgameController)
	route.AddTagRouter(router, oauthKey, controllers.TagController)
	route.AddCategoryRouter(router, oauthKey, controllers.CategoryController)
	route.AddMechanismRouter(router, oauthKey, controllers.MechanismController)

	log.Debug().Msg("setup complete")
	return &Base{
		router:      router,
		oauthHeader: oauthKey,
		dbMock:      mock,
	}
}
