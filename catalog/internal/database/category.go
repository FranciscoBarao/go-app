package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"

	"github.com/FranciscoBarao/catalog/internal/category"
	dbsql "github.com/FranciscoBarao/catalog/internal/database/sql"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/middleware"
)

// CreateCategory inserts a new category into the database.
func (p *Postgres) CreateCategory(ctx context.Context, c *category.Category) error {
	_, err := p.pool.Exec(ctx, dbsql.InsertCategory, c.Name)
	return mapPgError(err)
}

// GetCategory retrieves a single category by name.
func (p *Postgres) GetCategory(ctx context.Context, name string) (category.Category, error) {
	rows, err := p.pool.Query(ctx, dbsql.SelectCategory, name)
	if err != nil {
		return category.Category{}, mapPgError(err)
	}
	c, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[category.Category])
	return c, mapPgError(err)
}

// GetAllCategories retrieves all categories, optionally ordered by the given sort column.
func (p *Postgres) GetAllCategories(ctx context.Context, filter listopt.Params) ([]category.Category, error) {
	q := dbsql.SelectAllCategories
	var args []any

	if where, arg := filterClause(filter); where != "" {
		q += where
		args = append(args, arg)
	}

	if filter.Sort.Column != "" {
		q += fmt.Sprintf(dbsql.OrderBy, filter.Sort.Column, filter.Sort.Order)
	}

	logging.FromCtx(ctx).Debug().Str("query", q).Msg("GetAllCategories")

	rows, err := p.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, mapPgError(err)
	}
	categories, err := pgx.CollectRows(rows, pgx.RowToStructByPos[category.Category])
	return categories, mapPgError(err)
}

// DeleteCategory removes a category by name. Returns 404 if the category does not exist.
func (p *Postgres) DeleteCategory(ctx context.Context, name string) error {
	cmd, err := p.pool.Exec(ctx, dbsql.DeleteCategory, name)
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}
	return nil
}
