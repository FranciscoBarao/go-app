package main

import (
	"context"
	"net/http"

	_ "github.com/FranciscoBarao/catalog/build/docs"
	"github.com/FranciscoBarao/catalog/config"
	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/route"
	"github.com/FranciscoBarao/catalog/internal/store/postgres"
	"github.com/FranciscoBarao/catalog/internal/transport"
	"github.com/go-chi/chi/v5"
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
	cfg.MigrationPath = "internal/store/postgres/migrations"

	db, err := postgres.Connect(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer func() {
		db.Shutdown()
		log.Info().Msg("server stopped")
	}()

	oauthKey, port, err := config.OAuthAndPort()
	if err != nil {
		log.Fatal().Err(err).Msg("environment")
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
	router.Use(logging.RequestMiddleware())

	route.AddBoardGameRouter(router, oauthKey, bgController)
	route.AddCategoryRouter(router, oauthKey, categoryController)
	route.AddMechanismRouter(router, oauthKey, mechanismController)
	route.AddContributorRouter(router, oauthKey, contributorController)

	router.Get("/swagger/*", httpSwagger.Handler())

	log.Info().Str("port", port).Msg("server starting")
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal().Err(err).Msg("failed to create http server")
	}
}
