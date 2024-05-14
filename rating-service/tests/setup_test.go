package tests

import (
	"log"
	"os"

	"rating-service/controllers"
	"rating-service/database"
	"rating-service/repositories"
	"rating-service/route"
	"rating-service/services"

	"github.com/go-chi/chi/v5"
)

var router *chi.Mux
var oauthHeader string

// Prepares test environment
func init() {
	log.Println("Setup Starting")

	// Set Database
	db, err := database.Connect(false)
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
	route.AddRatingRouter(router, oauthKey, controllers.RatingController)

	log.Println("Setup Complete")
}
