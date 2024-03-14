package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"catalog/config"
	"catalog/controllers"
	"catalog/database"
	_ "catalog/docs"
	"catalog/repositories"
	"catalog/route"
	"catalog/services"
)

// @title Catalog App Swagger
// @version 1.0
// @description This microservice is a catalog for holding the possibly objects that can be used to create offers in the marketplace.

// @contact.name Francisco Barao
// @contact.email s.franciscobarao@gmail.com

// @BasePath /api/
func main() {
	// Fetch DB configs
	config, err := config.NewPostgresConfig()
	if err != nil {
		log.Fatal(err)
	}
	// Connect to Database
	db, err := database.Connect(config)
	if err != nil {
		log.Fatal("error occurred while connecting to database")
	}

	// Fetch Env variables
	oauthKey, oauthKeyPresent := os.LookupEnv("OAUTH_KEY")
	port, portPresent := os.LookupEnv("PORT")
	if !oauthKeyPresent || !portPresent {
		log.Fatal("error occurred while fetching essential env variables")
	}

	// Initialize Repositories & Services & controllers
	repositories := repositories.InitRepositories(db)
	services := services.InitServices(repositories)
	controllers := controllers.InitControllers(services)

	// Creates routing
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Adds Routers
	route.AddBoardGameRouter(router, oauthKey, controllers.BoardgameController)
	route.AddTagRouter(router, oauthKey, controllers.TagController)
	route.AddCategoryRouter(router, oauthKey, controllers.CategoryController)
	route.AddMechanismRouter(router, oauthKey, controllers.MechanismController)

	// documentation for developers
	router.Get("/swagger/*", httpSwagger.Handler())

	// Starts server
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("error occured while creating server: %s" + err.Error())

	}
	log.Println("server running on localhost:" + port)
}
