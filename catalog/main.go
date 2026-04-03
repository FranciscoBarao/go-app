package main

import (
	"context"
	"net/http"
	"os"

	logging "github.com/FranciscoBarao/catalog/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/FranciscoBarao/catalog/boardgame"
	"github.com/FranciscoBarao/catalog/category"
	"github.com/FranciscoBarao/catalog/config"
	"github.com/FranciscoBarao/catalog/database"
	_ "github.com/FranciscoBarao/catalog/docs"
	"github.com/FranciscoBarao/catalog/mechanism"
	"github.com/FranciscoBarao/catalog/route"
	"github.com/FranciscoBarao/catalog/tag"
	"github.com/FranciscoBarao/catalog/transport"
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
	db, err := database.Connect(config, &boardgame.Boardgame{}, &tag.Tag{}, &category.Category{}, &mechanism.Mechanism{}, &boardgame.Rating{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// Fetch Env variables
	oauthKey, oauthKeyPresent := os.LookupEnv("OAUTH_KEY")
	port, portPresent := os.LookupEnv("PORT")
	if !oauthKeyPresent || !portPresent {
		log.Fatal().Msg("failed to fetch essential env variables")
	}

	// Initialize Services
	tagSvc := tag.NewService(db)
	categorySvc := category.NewService(db)
	mechanismSvc := mechanism.NewService(db)
	boardgameSvc := boardgame.NewService(db, tagSvc, categorySvc, mechanismSvc)

	// Initialize Controllers
	bgController := transport.NewBoardgameController(boardgameSvc)
	tagController := transport.NewTagController(tagSvc)
	categoryController := transport.NewCategoryController(categorySvc)
	mechanismController := transport.NewMechanismController(mechanismSvc)

	// Creates routing
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Adds Routers
	route.AddBoardGameRouter(router, oauthKey, bgController)
	route.AddTagRouter(router, oauthKey, tagController)
	route.AddCategoryRouter(router, oauthKey, categoryController)
	route.AddMechanismRouter(router, oauthKey, mechanismController)

	// documentation for developers
	router.Get("/swagger/*", httpSwagger.Handler())

	// Starts server
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal().Err(err).Msg("failed to create http server")

	}
	log.Debug().Str("port", port).Str("ip", "localhost").Msg("server running")
}
