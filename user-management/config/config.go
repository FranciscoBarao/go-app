package config

import (
	"fmt"
	"os"
)

// LogLevel returns the configured log level, defaulting to "debug".
func LogLevel() string {
	if lvl, ok := os.LookupEnv("LOG_LEVEL"); ok {
		return lvl
	}
	return "debug"
}

// PostgresConfig holds the configuration for connecting to a PostgreSQL database.
type PostgresConfig struct {
	Host          string
	Username      string
	Password      string
	Port          string
	Database      string
	MigrationPath string
}

// NewPostgresConfig reads PostgreSQL connection settings from environment variables.
func NewPostgresConfig() (*PostgresConfig, error) {
	host, hostPresent := os.LookupEnv("DATABASE_HOST")
	user, userPresent := os.LookupEnv("POSTGRES_USER")
	pass, passPresent := os.LookupEnv("POSTGRES_PASSWORD")
	db, dbnamePresent := os.LookupEnv("POSTGRES_DB")
	port, portPresent := os.LookupEnv("DATABASE_PORT")

	if !hostPresent || !userPresent || !passPresent || !dbnamePresent || !portPresent {
		return nil, fmt.Errorf("failed to fetch postgres env vars")
	}

	return &PostgresConfig{
		Host:     host,
		Username: user,
		Password: pass,
		Database: db,
		Port:     port,
	}, nil
}
