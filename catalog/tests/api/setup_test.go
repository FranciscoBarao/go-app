package tests

import (
	"log"
	"testing"

	"catalog/controllers"
	"catalog/repositories"
	repositoriesMock "catalog/repositories/mock"
	"catalog/route"
	"catalog/services"

	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
)

const oauthKey = "secret-key"

type Base struct {
	router      *chi.Mux
	oauthHeader string
	dbMock      *repositoriesMock.MockDatabase
}

// Prepares test environment
func NewBase(t *testing.T) *Base {
	log.Println("Setup Starting")

	// Setup database mock
	mock := repositoriesMock.NewMockDatabase(gomock.NewController(t))

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

	log.Println("Setup Complete")
	return &Base{
		router:      router,
		oauthHeader: oauthKey,
		dbMock:      mock,
	}
}
