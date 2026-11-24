package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/oauth"

	"github.com/FranciscoBarao/user-management/config"
	"github.com/FranciscoBarao/user-management/internal/auth"
	"github.com/FranciscoBarao/user-management/internal/database"
	"github.com/FranciscoBarao/user-management/internal/logging"
	"github.com/FranciscoBarao/user-management/internal/route"
	"github.com/FranciscoBarao/user-management/internal/transport"
	"github.com/FranciscoBarao/user-management/internal/user"
)

func main() {
	logging.Init(config.LogLevel())

	ctx := context.Background()
	log := logging.FromCtx(ctx)

	// Fetch DB config
	cfg, err := config.NewPostgresConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to fetch database env variables")
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
	userSvc := user.NewService(db)

	// Initialize Auth Verifier
	verifier := auth.NewVerifier(userSvc)

	// Initialize Controllers
	userController := transport.NewUserController(userSvc)

	// OAuth server
	oauthServer := oauth.NewBearerServer(oauthKey, time.Minute*60, verifier, nil)

	// Creates routing
	router := chi.NewRouter()
	router.Use(chiMiddleware.Logger)

	// Auth routes
	router.Post("/api/login", oauthServer.UserCredentials)
	router.Post("/api/auth", oauthServer.ClientCredentials)

	// User routes
	route.AddUserRouter(router, oauthKey, userController)

	log.Debug().Msg("routes registered")

	// Starts server
	log.Info().Str("port", port).Msg("server starting")
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal().Err(err).Msg("failed to create http server")
	}
}
