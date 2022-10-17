package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"catalog/controllers"
	"catalog/database"
	_ "catalog/docs"
	"catalog/repositories"
	"catalog/route"
)

// @title Catalog App Swagger
// @version 1.0
// @description This microservice is a catalog for holding the possibly objects that can be used to create offers in the marketplace.

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
	repos := repositories.InitRepositories(db)
	controllers := controllers.InitControllers(repos)

	// Creates routing
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Adds Routers
	route.AddBoardGameRouter(router, controllers.BoardgameController)
	route.AddTagRouter(router, controllers.TagController)
	route.AddCategoryRouter(router, controllers.CategoryController)
	route.AddMechanismRouter(router, controllers.MechanismController)

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
