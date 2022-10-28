package main

import (
	"marketplace/controllers"
	"marketplace/database"
	"marketplace/repositories"
	"marketplace/route"
	"marketplace/services"

	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/oauth"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Marketplace App Swagger
// @version 1.0
// @description This microservice is a marketplace to create, display and buy offers

// @contact.name Francisco Barao
// @contact.email s.franciscobarao@gmail.com

// @BasePath /api/
func main() {
	// Connect to Database
	db, err := database.Connect()
	if err != nil {
		log.Println("Error occurred while connecting to database")
		return
	}

	// Initialize Repositories and controllers
	repositories := repositories.InitRepositories(db)
	services := services.InitServices(repositories)

	controllers := controllers.InitControllers(services)

	// Creates routing
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	// use the Bearer Authentication middleware
	router.Use(oauth.Authorize("mySecretKey-10101", nil))
	// Adds Routers
	route.AddOfferRouter(router, controllers.OfferController)

	// documentation for developers
	router.Get("/swagger/*", httpSwagger.Handler())

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
