package tests

import (
	"catalog/controllers"
	"catalog/database"
	"catalog/repositories"
	"catalog/route"
	"catalog/services"
	"log"

	"github.com/go-chi/chi/v5"
)

var router *chi.Mux

// Prepares test environment
func init() {
	log.Println("Setup Starting")

	// Set Database
	db, err := database.Connect()
	if err != nil {
		log.Println("Error occurred while connecting to database")
		return
	}

	// Set Repositories & Controllers & Services
	repositories := repositories.InitRepositories(db)
	services := services.InitServices(repositories)
	controllers := controllers.InitControllers(services)

	router = chi.NewRouter()
	// Adds Routers
	route.AddBoardGameRouter(router, controllers.BoardgameController)
	route.AddTagRouter(router, controllers.TagController)
	route.AddCategoryRouter(router, controllers.CategoryController)
	route.AddMechanismRouter(router, controllers.MechanismController)

	log.Println("Setup Complete")
}
