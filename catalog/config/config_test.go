package config

import (
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	full := map[string]string{
		"PORT":              "8080",
		"OAUTH_KEY":         "secret-key",
		"DATABASE_HOST":     "catalog-db",
		"POSTGRES_USER":     "user",
		"POSTGRES_PASSWORD": "password",
		"DATABASE_PORT":     "5432",
		"POSTGRES_DB":       "bg",
	}

	tests := []struct {
		name          string
		env           map[string]string
		wantLevel     string
		wantPort      int
		migrationPath string
		wantErr       []string
	}{
		{
			name:      "complete",
			env:       withEnv(full, "LOG_LEVEL", "info"),
			wantLevel: "info",
			wantPort:  8080,
		},
		{
			name:          "custom migration path",
			env:           withEnv(full, "MIGRATION_PATH", "db/migrations"),
			wantLevel:     "debug",
			wantPort:      8080,
			migrationPath: "db/migrations",
		},
		{
			name:      "default log level",
			env:       full,
			wantLevel: "debug",
			wantPort:  8080,
		},
		{
			name:      "empty log level",
			env:       withEnv(full, "LOG_LEVEL", ""),
			wantLevel: "debug",
			wantPort:  8080,
		},
		{
			name: "missing names",
			env:  map[string]string{"PORT": "8080"},
			wantErr: []string{
				"OAUTH_KEY",
				"DATABASE_HOST",
				"POSTGRES_USER",
				"POSTGRES_PASSWORD",
				"DATABASE_PORT",
				"POSTGRES_DB",
			},
		},
		{
			name:    "empty required value",
			env:     withEnv(full, "OAUTH_KEY", ""),
			wantErr: []string{"OAUTH_KEY"},
		},
		{
			name:    "port not a number",
			env:     withEnv(full, "PORT", "nope"),
			wantErr: []string{"Port"},
		},
		{
			name:    "port out of range",
			env:     withEnv(full, "PORT", "0"),
			wantErr: []string{"outside 1-65535"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := load(env.Options{Environment: tt.env})
			if len(tt.wantErr) > 0 {
				require.Error(t, err)
				for _, name := range tt.wantErr {
					assert.ErrorContains(t, err, name)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantLevel, cfg.LogLevel)
			assert.Equal(t, tt.wantPort, cfg.Port)
			assert.Equal(t, "catalog-db", cfg.Postgres.Host)
			assert.Equal(t, "user", cfg.Postgres.Username)
			assert.Equal(t, "password", cfg.Postgres.Password)
			assert.Equal(t, "5432", cfg.Postgres.Port)
			assert.Equal(t, "bg", cfg.Postgres.Database)
			wantPath := tt.migrationPath
			if wantPath == "" {
				wantPath = "internal/store/postgres/migrations"
			}
			assert.Equal(t, wantPath, cfg.Postgres.MigrationPath)
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr string
	}{
		{name: "lowest", port: 1},
		{name: "highest", port: 65535},
		{name: "zero", port: 0, wantErr: "outside 1-65535"},
		{name: "above range", port: 65536, wantErr: "outside 1-65535"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePort(tt.port)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func withEnv(base map[string]string, key string, value string) map[string]string {
	out := make(map[string]string, len(base)+1)
	for k, v := range base {
		out[k] = v
	}
	out[key] = value
	return out
}

func TestLoadReadsProcessEnvironment(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("OAUTH_KEY", "secret-key")
	t.Setenv("DATABASE_HOST", "catalog-db")
	t.Setenv("POSTGRES_USER", "user")
	t.Setenv("POSTGRES_PASSWORD", "password")
	t.Setenv("DATABASE_PORT", "5432")
	t.Setenv("POSTGRES_DB", "bg")
	t.Setenv("LOG_LEVEL", "warn")

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "warn", cfg.LogLevel)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "secret-key", cfg.OAuthKey)
}
