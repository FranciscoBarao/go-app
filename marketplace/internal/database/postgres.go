package database

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/FranciscoBarao/marketplace/config"
	"github.com/FranciscoBarao/marketplace/internal/logging"
	"github.com/FranciscoBarao/marketplace/internal/middleware"
)

const (
	uniqueViolationCode     = "23505"
	foreignKeyViolationCode = "23503"
)

// Postgres holds the pgx connection pool.
type Postgres struct {
	pool *pgxpool.Pool
}

// Connect creates a new Postgres instance with an initialized connection pool.
func Connect(ctx context.Context, cfg *config.PostgresConfig) (*Postgres, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database,
	)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("database connect: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping: %w", err)
	}

	if cfg.MigrationPath != "" {
		if err := runMigrations(ctx, pool, cfg.MigrationPath); err != nil {
			pool.Close()
			return nil, fmt.Errorf("database migrate: %w", err)
		}
	}

	logging.FromCtx(ctx).Debug().Str("host", cfg.Host).Str("port", cfg.Port).Str("database", cfg.Database).Msg("database connected")
	return &Postgres{pool: pool}, nil
}

// Close shuts down the connection pool.
func (p *Postgres) Close() {
	p.pool.Close()
}

// runMigrations reads all *.up.sql files from dir in sorted order and executes them.
func runMigrations(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migration dir: %w", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, f := range files {
		sql, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f, err)
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("exec migration %s: %w", f, err)
		}
		logging.FromCtx(ctx).Debug().Str("file", f).Msg("migration applied")
	}

	return nil
}

// mapPgError converts pgx errors to domain errors.
func mapPgError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		logging.FromCtx(context.Background()).Error().Str("code", pgErr.Code).Str("message", pgErr.Message).Msg("postgres error")
		switch pgErr.Code {
		case uniqueViolationCode:
			return middleware.NewError(http.StatusConflict, "entry already registered")
		case foreignKeyViolationCode:
			return middleware.NewError(http.StatusConflict, "referenced record not found")
		}
	}
	return err
}
