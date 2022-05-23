package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/unrolled/render"

	"go-app/controller"
	"go-app/database"
	_ "go-app/docs"
	"go-app/repository"
	"go-app/route"
)

// @title Boardgame App Swagger
// @version 1.0
// @description This is adescription

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
	repos := repository.InitRepositories(db)
	controllers := controller.InitControllers(repos)

	// Creates routing
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Adds Routers
	route.AddBoardGameRouter(router, controllers.BoardgameController)
	route.AddTagRouter(router, controllers.TagController)

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

func Test(w http.ResponseWriter, r *http.Request) {
	render.New().JSON(w, http.StatusOK, "200")
}
