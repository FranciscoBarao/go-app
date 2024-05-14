package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/FranciscoBarao/catalog/config"
	"github.com/FranciscoBarao/catalog/controllers"
	"github.com/FranciscoBarao/catalog/database"
	_ "github.com/FranciscoBarao/catalog/docs"
	logging "github.com/FranciscoBarao/catalog/middleware/logging"
	"github.com/FranciscoBarao/catalog/repositories"
	"github.com/FranciscoBarao/catalog/route"
	"github.com/FranciscoBarao/catalog/services"
)

// @title Catalog App Swagger
// @version 1.0
// @description This microservice is a catalog for holding the possibly objects that can be used to create offers in the marketplace.
// @contact.name Francisco Barao
// @contact.email s.franciscobarao@gmail.com
// @BasePath /api/
func main() {
	ctx := context.Background()
	log := logging.FromCtx(ctx)
	// Fetch DB configs
	config, err := config.NewPostgresConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to fetch database env variables")
	}
	// Connect to Database
	db, err := database.Connect(config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// Fetch Env variables
	oauthKey, oauthKeyPresent := os.LookupEnv("OAUTH_KEY")
	port, portPresent := os.LookupEnv("PORT")
	if !oauthKeyPresent || !portPresent {
		log.Fatal().Msg("failed to fetch essential env variables")
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
		log.Fatal().Err(err).Msg("failed to create http server")

	}
	log.Debug().Str("port", port).Str("ip", "localhost").Msg("server running")
}
