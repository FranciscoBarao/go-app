package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"

	"github.com/FranciscoBarao/catalog/internal/category"
	dbsql "github.com/FranciscoBarao/catalog/internal/database/sql"
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
func (p *Postgres) GetAllCategories(ctx context.Context, sort string) ([]category.Category, error) {
	col, err := categorySortColumn(sort)
	if err != nil {
		return nil, err
	}

	q := dbsql.SelectAllCategories
	if col != "" {
		q += fmt.Sprintf(dbsql.OrderBy, col)
	}

	rows, err := p.pool.Query(ctx, q)
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

// categorySortColumn validates and returns the sort column for categories.
func categorySortColumn(sort string) (string, error) {
	if sort == "" {
		return "", nil
	}
	switch sort {
	case "name", "created_at", "updated_at":
		return sort, nil
	default:
		return "", middleware.NewError(http.StatusBadRequest, fmt.Sprintf("invalid sort column: %s", sort))
	}
}
