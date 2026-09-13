package database

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/FranciscoBarao/catalog/internal/contributor"
	dbsql "github.com/FranciscoBarao/catalog/internal/database/sql"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/middleware"
)

// CreateContributor inserts a new contributor.
func (p *Postgres) CreateContributor(ctx context.Context, c *contributor.Contributor) error {
	err := p.pool.QueryRow(ctx, dbsql.InsertContributor, c.Slug, c.Name, c.Bio, c.BggID).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	return mapPgError(err)
}

// UpdateContributor updates contributor fields by slug and refreshes updated_at.
func (p *Postgres) UpdateContributor(ctx context.Context, c *contributor.Contributor) error {
	err := p.pool.QueryRow(ctx, dbsql.UpdateContributor, c.Name, c.Bio, c.BggID, c.Slug).
		Scan(&c.UpdatedAt)
	return mapPgError(err)
}

// GetContributorBySlug retrieves a contributor by slug.
func (p *Postgres) GetContributorBySlug(ctx context.Context, slug string) (contributor.Contributor, error) {
	rows, err := p.pool.Query(ctx, dbsql.SelectContributorBySlug, slug)
	if err != nil {
		return contributor.Contributor{}, mapPgError(err)
	}
	c, err := pgx.CollectOneRow(rows, scanContributorRow)
	return c, mapPgError(err)
}

// GetContributorIDBySlug retrieves only the id of an active contributor by slug.
func (p *Postgres) GetContributorIDBySlug(ctx context.Context, slug string) (uint, error) {
	var id uint
	err := p.pool.QueryRow(ctx, dbsql.SelectContributorIDBySlug, slug).Scan(&id)
	return id, mapPgError(err)
}

// GetAllContributors retrieves active contributors with optional filtering,
// sorting, and pagination, returning the page and the total count of matching rows.
func (p *Postgres) GetAllContributors(ctx context.Context, query listopt.Query) ([]contributor.Contributor, int, error) {
	countQuery, countArgs := buildCountQuery(dbsql.CountContributors, query)
	selectQuery, selectArgs := buildPaginatedQuery(dbsql.SelectAllContributors, query)

	logging.FromCtx(ctx).Debug().Str("query", selectQuery).Msg("GetAllContributors")

	var total int
	if err := p.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, mapPgError(err)
	}

	rows, err := p.pool.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, mapPgError(err)
	}
	contributors, err := pgx.CollectRows(rows, scanContributorRow)
	return contributors, total, mapPgError(err)
}

// DeleteContributor soft-deletes or hard-deletes a contributor.
func (p *Postgres) DeleteContributor(ctx context.Context, slug string, hard bool) error {
	var (
		cmd pgconn.CommandTag
		err error
	)
	if hard {
		cmd, err = p.pool.Exec(ctx, dbsql.HardDeleteContributor, slug)
	} else {
		cmd, err = p.pool.Exec(ctx, dbsql.SoftDeleteContributor, slug)
	}
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}
	return nil
}

func scanContributorRow(row pgx.CollectableRow) (contributor.Contributor, error) {
	var c contributor.Contributor
	err := row.Scan(&c.ID, &c.Slug, &c.Name, &c.Bio, &c.BggID, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
	return c, err
}
