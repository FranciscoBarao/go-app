package config

import (
	"fmt"
	"os"
	"strings"
)

// LogLevel returns LOG_LEVEL. An unset value is debug, for local development.
// A deployed environment sets LOG_LEVEL=info.
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

	var missing []string
	missing = appendMissing(missing, "DATABASE_HOST", hostPresent)
	missing = appendMissing(missing, "POSTGRES_USER", userPresent)
	missing = appendMissing(missing, "POSTGRES_PASSWORD", passPresent)
	missing = appendMissing(missing, "POSTGRES_DB", dbnamePresent)
	missing = appendMissing(missing, "DATABASE_PORT", portPresent)
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}

	return &PostgresConfig{
		Host:     host,
		Username: user,
		Password: pass,
		Database: db,
		Port:     port,
	}, nil
}

// OAuthAndPort reads OAUTH_KEY and PORT. A missing variable is named in the error.
func OAuthAndPort() (string, string, error) {
	oauthKey, oauthKeyPresent := os.LookupEnv("OAUTH_KEY")
	port, portPresent := os.LookupEnv("PORT")

	var missing []string
	missing = appendMissing(missing, "OAUTH_KEY", oauthKeyPresent)
	missing = appendMissing(missing, "PORT", portPresent)
	if len(missing) > 0 {
		return "", "", fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}
	return oauthKey, port, nil
}

func appendMissing(missing []string, name string, present bool) []string {
	if present {
		return missing
	}
	return append(missing, name)
}

func (p *PostgresConfig) String() string {
	return "host=" + p.Host + " user=" + p.Username + " password=" + p.Password + " dbname=" + p.Database + " port=" + p.Port
}
