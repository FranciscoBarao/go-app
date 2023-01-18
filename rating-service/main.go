package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpSwagger "github.com/swaggo/http-swagger"

	"rating-service/controllers"
	"rating-service/database"
	_ "rating-service/docs"
	"rating-service/repositories"
	"rating-service/route"
	"rating-service/services"
)

// @title Rating Service App Swagger
// @version 1.0
// @description This microservice is an abstracted way of rating other services' products in the architecture.

// @contact.name Francisco Barao
// @contact.email s.franciscobarao@gmail.com

// @BasePath /api/
func main() {
	// Connect to Database
	db, err := database.Connect(false)
	if err != nil {
		log.Println("Error occurred while connecting to database")
		return
	}

	// Fetch Env variables
	oauthKey, oauthKeyPresent := os.LookupEnv("OAUTH_KEY")
	port, portPresent := os.LookupEnv("PORT")
	if !oauthKeyPresent || !portPresent {
		log.Println("Error occurred while fetching essential env variables")
		return
	}

	// Initialize Repositories & Services & controllers
	repositories := repositories.InitRepositories(db)
	services := services.InitServices(repositories)
	controllers := controllers.InitControllers(services)

	// Creates routing
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Adds Routers
	route.AddRatingRouter(router, oauthKey, controllers.RatingController)

	// documentation for developers
	router.Get("/swagger/*", httpSwagger.Handler())

	// Starts server
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Println("Error occured while creating Server" + err.Error())
		return
	}
	log.Println("Server is Running on localhost:" + port)
}
