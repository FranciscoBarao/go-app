package main

import (
	"context"
	"net/http"
	"os"

	"github.com/FranciscoBarao/catalog/config"
	_ "github.com/FranciscoBarao/catalog/docs"
	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/database"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/route"
	"github.com/FranciscoBarao/catalog/internal/transport"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Catalog App Swagger
// @version 1.0
// @description Catalog service for boardgames and related metadata.
// @contact.name Francisco Barao
// @contact.email s.franciscobarao@gmail.com
// @BasePath /api/
func main() {
	logging.Init(config.LogLevel())

	ctx := context.Background()
	log := logging.FromCtx(ctx)

	cfg, err := config.NewPostgresConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to fetch database env variables")
	}
	cfg.MigrationPath = "internal/database/migrations"

	db, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	oauthKey, oauthKeyPresent := os.LookupEnv("OAUTH_KEY")
	port, portPresent := os.LookupEnv("PORT")
	if !oauthKeyPresent || !portPresent {
		log.Fatal().Msg("failed to fetch essential env variables")
	}

	categorySvc := category.NewService(db)
	mechanismSvc := mechanism.NewService(db)
	contributorSvc := contributor.NewService(db)
	boardgameSvc := boardgame.NewService(db, categorySvc, mechanismSvc, contributorSvc)

	bgController := transport.NewBoardgameController(boardgameSvc)
	categoryController := transport.NewCategoryController(categorySvc)
	mechanismController := transport.NewMechanismController(mechanismSvc)
	contributorController := transport.NewContributorController(contributorSvc)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	route.AddBoardGameRouter(router, oauthKey, bgController)
	route.AddCategoryRouter(router, oauthKey, categoryController)
	route.AddMechanismRouter(router, oauthKey, mechanismController)
	route.AddContributorRouter(router, oauthKey, contributorController)

	router.Get("/swagger/*", httpSwagger.Handler())

	log.Debug().Msg("routes registered")
	log.Info().Str("port", port).Msg("server starting")
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal().Err(err).Msg("failed to create http server")
	}
}
