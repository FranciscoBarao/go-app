package postgres

import (
	"context"
	"github.com/FranciscoBarao/catalog/internal/store/postgres/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/store"
)

// CreateMechanism inserts a new mechanism into the database.
func (p *Postgres) CreateMechanism(ctx context.Context, m *mechanism.Mechanism) error {
	err := p.pool.QueryRow(ctx, sql.InsertMechanism, m.Slug, m.Name, m.BggID).
		Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	return mapPgError(ctx, err)
}

// GetMechanismBySlug retrieves a single mechanism by slug.
func (p *Postgres) GetMechanismBySlug(ctx context.Context, slug string) (mechanism.Mechanism, error) {
	rows, err := p.pool.Query(ctx, sql.SelectMechanismBySlug, slug)
	if err != nil {
		return mechanism.Mechanism{}, mapPgError(ctx, err)
	}
	m, err := pgx.CollectOneRow(rows, scanMechanismRow)
	return m, mapPgError(ctx, err)
}

// GetMechanismIDBySlug retrieves only the id of an active mechanism by slug.
func (p *Postgres) GetMechanismIDBySlug(ctx context.Context, slug string) (uint, error) {
	var id uint
	err := p.pool.QueryRow(ctx, sql.SelectMechanismIDBySlug, slug).Scan(&id)
	return id, mapPgError(ctx, err)
}

// GetAllMechanisms retrieves active mechanisms with optional filtering, sorting,
// and pagination, returning the page and the total count of matching rows.
func (p *Postgres) GetAllMechanisms(ctx context.Context, query listopt.Query) ([]mechanism.Mechanism, int, error) {
	countQuery, countArgs := buildCountQuery(sql.CountMechanisms, query)
	selectQuery, selectArgs := buildPaginatedQuery(sql.SelectAllMechanisms, query)

	var total int
	if err := p.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, mapPgError(ctx, err)
	}

	rows, err := p.pool.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, mapPgError(ctx, err)
	}
	mechanisms, err := pgx.CollectRows(rows, scanMechanismRow)
	return mechanisms, total, mapPgError(ctx, err)
}

// DeleteMechanism soft-deletes or hard-deletes a mechanism by slug.
func (p *Postgres) DeleteMechanism(ctx context.Context, slug string, hard bool) error {
	var (
		cmd pgconn.CommandTag
		err error
	)
	if hard {
		cmd, err = p.pool.Exec(ctx, sql.HardDeleteMechanism, slug)
	} else {
		cmd, err = p.pool.Exec(ctx, sql.SoftDeleteMechanism, slug)
	}
	if err != nil {
		return mapPgError(ctx, err)
	}
	if cmd.RowsAffected() == 0 {
		return store.ErrNotFound
	}
	return nil
}

func scanMechanismRow(row pgx.CollectableRow) (mechanism.Mechanism, error) {
	var m mechanism.Mechanism
	err := row.Scan(&m.ID, &m.Slug, &m.Name, &m.BggID, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt)
	return m, err
}
