package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/render"

	"go-app/controller"
	"go-app/database"
	"go-app/repository"
)

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
	router.Post("/api/boardgame", controllers.BoardgameController.CreateBG)
	router.Get("/api/boardgame", controllers.BoardgameController.GetBG)
	router.Patch("/api/boardgame", controllers.BoardgameController.UpdateBG)
	router.Delete("/api/boardgame", controllers.BoardgameController.DeleteBG)

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
