package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"

	dbsql "github.com/FranciscoBarao/catalog/internal/database/sql"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/middleware"
)

// CreateMechanism inserts a new mechanism into the database.
func (p *Postgres) CreateMechanism(ctx context.Context, m *mechanism.Mechanism) error {
	_, err := p.pool.Exec(ctx, dbsql.InsertMechanism, m.Name)
	return mapPgError(err)
}

// GetMechanism retrieves a single mechanism by name.
func (p *Postgres) GetMechanism(ctx context.Context, name string) (mechanism.Mechanism, error) {
	rows, err := p.pool.Query(ctx, dbsql.SelectMechanism, name)
	if err != nil {
		return mechanism.Mechanism{}, mapPgError(err)
	}
	m, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[mechanism.Mechanism])
	return m, mapPgError(err)
}

// GetAllMechanisms retrieves all mechanisms, optionally ordered by the given sort column.
func (p *Postgres) GetAllMechanisms(ctx context.Context, filter listopt.Params) ([]mechanism.Mechanism, error) {

	q := dbsql.SelectAllMechanisms
	var args []any

	if where, arg := filterClause(filter); where != "" {
		q += where
		args = append(args, arg)
	}

	if filter.Sort.Column != "" {
		q += fmt.Sprintf(dbsql.OrderBy, filter.Sort.Column, filter.Sort.Order)
	}

	rows, err := p.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, mapPgError(err)
	}
	mechanisms, err := pgx.CollectRows(rows, pgx.RowToStructByPos[mechanism.Mechanism])
	return mechanisms, mapPgError(err)
}

// DeleteMechanism removes a mechanism by name. Returns 404 if the mechanism does not exist.
func (p *Postgres) DeleteMechanism(ctx context.Context, name string) error {
	cmd, err := p.pool.Exec(ctx, dbsql.DeleteMechanism, name)
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}
	return nil
}
