package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/FranciscoBarao/marketplace/config"
	_ "github.com/FranciscoBarao/marketplace/docs"
	"github.com/FranciscoBarao/marketplace/internal/database"
	"github.com/FranciscoBarao/marketplace/internal/logging"
	"github.com/FranciscoBarao/marketplace/internal/offer"
	"github.com/FranciscoBarao/marketplace/internal/route"
	"github.com/FranciscoBarao/marketplace/internal/transport"
)

// @title Marketplace App Swagger
// @version 1.0
// @description This microservice is a marketplace to create, display and buy offers

// @contact.name Francisco Barao
// @contact.email s.franciscobarao@gmail.com

// @BasePath /api/
func main() {
	logging.Init(config.LogLevel())

	ctx := context.Background()
	log := logging.FromCtx(ctx)

	// Fetch DB config
	cfg, err := config.NewPostgresConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to configure database")
	}
	cfg.MigrationPath = "database/migrations"

	// Connect to Database
	db, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	// Fetch Env variables
	oauthKey, oauthKeyPresent := os.LookupEnv("OAUTH_KEY")
	port, portPresent := os.LookupEnv("PORT")
	if !oauthKeyPresent || !portPresent {
		log.Fatal().Msg("failed to fetch essential env variables")
	}

	// Initialize Services
	offerSvc := offer.NewService(db)

	// Initialize Controllers
	offerController := transport.NewOfferController(offerSvc)

	// Creates routing
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Adds Routers
	route.AddOfferRouter(router, oauthKey, offerController)

	// documentation for developers
	router.Get("/swagger/*", httpSwagger.Handler())

	log.Debug().Msg("routes registered")

	// Starts server
	log.Info().Str("port", port).Msg("server starting")
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal().Err(err).Msg("failed to create http server")
	}
}
