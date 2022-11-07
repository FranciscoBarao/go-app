package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"user-management/controllers"
	"user-management/database"

	"user-management/repositories"
	"user-management/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/oauth"
)

func main() {
	// Connect to Database
	db, err := database.Connect()
	if err != nil {
		log.Println("Error occurred while connecting to database")
		return
	}

	// Fetch oauth key
	oauthKey, oauthKeyPresent := os.LookupEnv("OAUTH_KEY")
	if !oauthKeyPresent {
		log.Println("Error occurred while fetching Oauth Key")
		return
	}

	// Initialize Repositories & Services & controllers
	repositories := repositories.InitRepositories(db)
	services := services.InitServices(repositories)
	controllers := controllers.InitControllers(services)

	// Creates routing
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	oauthServer := oauth.NewBearerServer(
		oauthKey,
		time.Minute*60,
		&controllers.VerifierController,
		nil)

	// Adds Routers
	router.Post("/api/register", controllers.UserController.Register)
	router.Post("/api/login", oauthServer.UserCredentials)
	router.Post("/api/auth", oauthServer.ClientCredentials)

	// Starts server
	port, portPresent := os.LookupEnv("PORT")
	if !portPresent {
		log.Println("Error occurred while fetching Port")
		return
	}

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Println("Error occured while creating Server" + err.Error())
		return
	}
	log.Println("Server is Running on localhost:" + port)
}
