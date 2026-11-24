package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config is the process configuration, read once from the environment at startup.
type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
	Port     int    `env:"PORT,required"`
	OAuthKey string `env:"OAUTH_KEY,required,notEmpty"`
	Postgres PostgresConfig
}

// PostgresConfig is the PostgreSQL connection configuration.
type PostgresConfig struct {
	Host          string `env:"DATABASE_HOST,required,notEmpty"`
	Username      string `env:"POSTGRES_USER,required,notEmpty"`
	Password      string `env:"POSTGRES_PASSWORD,required,notEmpty"`
	Port          string `env:"DATABASE_PORT,required,notEmpty"`
	Database      string `env:"POSTGRES_DB,required,notEmpty"`
	MigrationPath string `env:"MIGRATION_PATH" envDefault:"internal/store/postgres/migrations"`
}

// Load reads the process environment into Config.
// Every missing required variable and every value that fails to parse is named in the returned error.
// An unset or empty LOG_LEVEL is debug. An unknown level is left as read; logging.SetLevel warns and uses debug.
func Load() (*Config, error) {
	return load(env.Options{})
}

func load(opts env.Options) (*Config, error) {
	var cfg Config
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		return nil, err
	}
	if err := validatePort(cfg.Port); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port %d is outside 1-65535", port)
	}
	return nil
}
