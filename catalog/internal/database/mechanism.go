package database

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	dbsql "github.com/FranciscoBarao/catalog/internal/database/sql"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/middleware"
)

// CreateMechanism inserts a new mechanism into the database.
func (p *Postgres) CreateMechanism(ctx context.Context, m *mechanism.Mechanism) error {
	err := p.pool.QueryRow(ctx, dbsql.InsertMechanism, m.Slug, m.Name, m.BggID).
		Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	return mapPgError(err)
}

// GetMechanismBySlug retrieves a single mechanism by slug.
func (p *Postgres) GetMechanismBySlug(ctx context.Context, slug string) (mechanism.Mechanism, error) {
	rows, err := p.pool.Query(ctx, dbsql.SelectMechanismBySlug, slug)
	if err != nil {
		return mechanism.Mechanism{}, mapPgError(err)
	}
	m, err := pgx.CollectOneRow(rows, scanMechanismRow)
	return m, mapPgError(err)
}

// GetMechanismIDBySlug retrieves only the id of an active mechanism by slug.
func (p *Postgres) GetMechanismIDBySlug(ctx context.Context, slug string) (uint, error) {
	var id uint
	err := p.pool.QueryRow(ctx, dbsql.SelectMechanismIDBySlug, slug).Scan(&id)
	return id, mapPgError(err)
}

// GetAllMechanisms retrieves all mechanisms.
func (p *Postgres) GetAllMechanisms(ctx context.Context, filter listopt.Params) ([]mechanism.Mechanism, error) {
	q := dbsql.SelectAllMechanisms
	var args []any

	if where, arg := filterClause(filter); where != "" {
		q += strings.Replace(where, " WHERE ", " AND ", 1)
		args = append(args, arg)
	}

	if filter.Sort.Column != "" {
		q += fmt.Sprintf(dbsql.OrderBy, filter.Sort.Column, filter.Sort.Order)
	}

	logging.FromCtx(ctx).Debug().Str("query", q).Msg("GetAllMechanisms")

	rows, err := p.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, mapPgError(err)
	}
	mechanisms, err := pgx.CollectRows(rows, scanMechanismRow)
	return mechanisms, mapPgError(err)
}

// DeleteMechanism soft-deletes or hard-deletes a mechanism by slug.
func (p *Postgres) DeleteMechanism(ctx context.Context, slug string, hard bool) error {
	var (
		cmd pgconn.CommandTag
		err error
	)
	if hard {
		cmd, err = p.pool.Exec(ctx, dbsql.HardDeleteMechanism, slug)
	} else {
		cmd, err = p.pool.Exec(ctx, dbsql.SoftDeleteMechanism, slug)
	}
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}
	return nil
}

func scanMechanismRow(row pgx.CollectableRow) (mechanism.Mechanism, error) {
	var m mechanism.Mechanism
	err := row.Scan(&m.ID, &m.Slug, &m.Name, &m.BggID, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt)
	return m, err
}
