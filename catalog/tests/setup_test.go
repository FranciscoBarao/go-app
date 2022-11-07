package tests

import (
	"catalog/controllers"
	"catalog/database"
	"catalog/repositories"
	"catalog/route"
	"catalog/services"
	"log"
	"os"

	"github.com/go-chi/chi/v5"
)

var router *chi.Mux
var oauthHeader string

// Prepares test environment
func init() {
	log.Println("Setup Starting")

	// Set Database
	db, err := database.Connect()
	if err != nil {
		log.Println("Error occurred while connecting to database")
		return
	}

	// Fetch Oauth Key
	oauthKey, oauthKeyPresent := os.LookupEnv("OAUTH_KEY")
	header, oauthHeaderPresent := os.LookupEnv("OAUTH_HEADER_TEST")
	if !oauthKeyPresent || !oauthHeaderPresent {
		log.Println("Error occurred while fetching essential env variables")
		return
	}

	// Defines Token header to be used in requests
	oauthHeader = header

	// Set Repositories & Controllers & Services
	repositories := repositories.InitRepositories(db)
	services := services.InitServices(repositories)
	controllers := controllers.InitControllers(services)

	router = chi.NewRouter()

	// Adds Routers
	route.AddBoardGameRouter(router, oauthKey, controllers.BoardgameController)
	route.AddTagRouter(router, oauthKey, controllers.TagController)
	route.AddCategoryRouter(router, oauthKey, controllers.CategoryController)
	route.AddMechanismRouter(router, oauthKey, controllers.MechanismController)

	log.Println("Setup Complete")
}
