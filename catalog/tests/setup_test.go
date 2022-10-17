package tests

import (
	"catalog/controllers"
	"catalog/database"
	"catalog/repositories"
	"catalog/route"
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

	// Set Repositories & Controllers
	repositories := repositories.InitRepositories(db)
	controllers := controllers.InitControllers(repositories)

	router = chi.NewRouter()
	// Adds Routers
	route.AddBoardGameRouter(router, controllers.BoardgameController)
	route.AddTagRouter(router, controllers.TagController)
	route.AddCategoryRouter(router, controllers.CategoryController)
	route.AddMechanismRouter(router, controllers.MechanismController)

	log.Println("Setup Complete")
}
