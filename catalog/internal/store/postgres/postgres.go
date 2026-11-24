package postgres

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // register the pgx5 migrate driver
	_ "github.com/golang-migrate/migrate/v4/source/file"     // register the file migrate source
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/FranciscoBarao/catalog/config"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/store"
)

const (
	uniqueViolationCode     = "23505"
	foreignKeyViolationCode = "23503"
)

// database is the subset of a connection pool used by store methods.
type database interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Postgres is the PostgreSQL store.
type Postgres struct {
	pool  database
	close func()
}

// Connect opens a pool and applies pending up migrations.
// No pending migrations is a successful result.
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
		if err := applyMigrations(cfg); err != nil {
			pool.Close()
			return nil, fmt.Errorf("database migrate: %w", err)
		}
	}

	logging.FromCtx(ctx).Info().Str("host", cfg.Host).Str("port", cfg.Port).Str("database", cfg.Database).Msg("database connected")
	return &Postgres{pool: pool, close: pool.Close}, nil
}

// Shutdown closes the connection pool.
func (p *Postgres) Shutdown() {
	if p.close != nil {
		p.close()
	}
}

func applyMigrations(cfg *config.PostgresConfig) error {
	abs, err := filepath.Abs(cfg.MigrationPath)
	if err != nil {
		return fmt.Errorf("migration path: %w", err)
	}

	migrator, err := migrate.New("file://"+filepath.ToSlash(abs), databaseURL(cfg))
	if err != nil {
		return fmt.Errorf("open migrations: %w", err)
	}
	defer func() { _, _ = migrator.Close() }()

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

func databaseURL(cfg *config.PostgresConfig) string {
	u := &url.URL{
		Scheme:   "pgx5",
		User:     url.UserPassword(cfg.Username, cfg.Password),
		Host:     net.JoinHostPort(cfg.Host, cfg.Port),
		Path:     "/" + cfg.Database,
		RawQuery: "sslmode=disable",
	}
	return u.String()
}

// mapPgError converts pgx errors to store errors.
// Unexpected driver errors are logged and replaced so they do not leave the adapter.
func mapPgError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return store.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case uniqueViolationCode:
			return store.ErrAlreadyExists
		case foreignKeyViolationCode:
			return store.ErrForeignKey
		}
		logging.FromCtx(ctx).Error().Str("code", pgErr.Code).Str("message", pgErr.Message).Msg("postgres error")
		return errors.New("database operation failed")
	}
	logging.FromCtx(ctx).Error().Err(err).Msg("database operation failed")
	return errors.New("database operation failed")
}
