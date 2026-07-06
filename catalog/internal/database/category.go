package database

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/FranciscoBarao/catalog/internal/category"
	dbsql "github.com/FranciscoBarao/catalog/internal/database/sql"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/middleware"
)

// CreateCategory inserts a new category into the database.
func (p *Postgres) CreateCategory(ctx context.Context, c *category.Category) error {
	err := p.pool.QueryRow(ctx, dbsql.InsertCategory, c.Slug, c.Name, c.BggID).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	return mapPgError(err)
}

// GetCategoryBySlug retrieves a single category by slug.
func (p *Postgres) GetCategoryBySlug(ctx context.Context, slug string) (category.Category, error) {
	rows, err := p.pool.Query(ctx, dbsql.SelectCategoryBySlug, slug)
	if err != nil {
		return category.Category{}, mapPgError(err)
	}
	c, err := pgx.CollectOneRow(rows, scanCategoryRow)
	return c, mapPgError(err)
}

// GetCategoryIDBySlug retrieves only the id of an active category by slug.
func (p *Postgres) GetCategoryIDBySlug(ctx context.Context, slug string) (uint, error) {
	var id uint
	err := p.pool.QueryRow(ctx, dbsql.SelectCategoryIDBySlug, slug).Scan(&id)
	return id, mapPgError(err)
}

// GetAllCategories retrieves all categories, optionally filtered and sorted.
func (p *Postgres) GetAllCategories(ctx context.Context, filter listopt.Params) ([]category.Category, error) {
	q := dbsql.SelectAllCategories
	var args []any

	if where, arg := filterClause(filter); where != "" {
		q += strings.Replace(where, " WHERE ", " AND ", 1)
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
	categories, err := pgx.CollectRows(rows, scanCategoryRow)
	return categories, mapPgError(err)
}

// DeleteCategory soft-deletes or hard-deletes a category by slug.
func (p *Postgres) DeleteCategory(ctx context.Context, slug string, hard bool) error {
	var (
		cmd pgconn.CommandTag
		err error
	)
	if hard {
		cmd, err = p.pool.Exec(ctx, dbsql.HardDeleteCategory, slug)
	} else {
		cmd, err = p.pool.Exec(ctx, dbsql.SoftDeleteCategory, slug)
	}
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}
	return nil
}

func scanCategoryRow(row pgx.CollectableRow) (category.Category, error) {
	var c category.Category
	err := row.Scan(&c.ID, &c.Slug, &c.Name, &c.BggID, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
	return c, err
}
