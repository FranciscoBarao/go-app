package main

import (
	"context"
	"net/http"
	"strconv"

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
	logging.Init()
	ctx := context.Background()
	log := logging.FromCtx(ctx)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	logging.SetLevel(cfg.LogLevel)

	db, err := postgres.Connect(ctx, &cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer func() {
		db.Shutdown()
		log.Info().Msg("server stopped")
	}()

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

	route.AddBoardGameRouter(router, cfg.OAuthKey, bgController)
	route.AddCategoryRouter(router, cfg.OAuthKey, categoryController)
	route.AddMechanismRouter(router, cfg.OAuthKey, mechanismController)
	route.AddContributorRouter(router, cfg.OAuthKey, contributorController)

	router.Get("/swagger/*", httpSwagger.Handler())

	log.Info().Int("port", cfg.Port).Msg("server starting")
	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), router); err != nil {
		log.Fatal().Err(err).Msg("failed to create http server")
	}
}
