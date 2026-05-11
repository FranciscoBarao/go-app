package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"

	dbsql "github.com/FranciscoBarao/catalog/internal/database/sql"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/query"
	"github.com/FranciscoBarao/catalog/internal/tag"
)

// CreateTag inserts a new tag into the database.
func (p *Postgres) CreateTag(ctx context.Context, t *tag.Tag) error {
	_, err := p.pool.Exec(ctx, dbsql.InsertTag, t.Name)
	return mapPgError(err)
}

// GetTag retrieves a single tag by name.
func (p *Postgres) GetTag(ctx context.Context, name string) (tag.Tag, error) {
	rows, err := p.pool.Query(ctx, dbsql.SelectTag, name)
	if err != nil {
		return tag.Tag{}, mapPgError(err)
	}
	t, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[tag.Tag])
	return t, mapPgError(err)
}

// GetAllTags retrieves all tags, optionally ordered by the given sort column.
func (p *Postgres) GetAllTags(ctx context.Context, filter query.Filter) ([]tag.Tag, error) {
	col, err := tagSortColumn(filter.SortColumn)
	if err != nil {
		return nil, err
	}

	q := dbsql.SelectAllTags
	if col != "" {
		q += fmt.Sprintf(dbsql.OrderBy, col, filter.SortOrder)
	}

	rows, err := p.pool.Query(ctx, q)
	if err != nil {
		return nil, mapPgError(err)
	}
	tags, err := pgx.CollectRows(rows, pgx.RowToStructByPos[tag.Tag])
	return tags, mapPgError(err)
}

// DeleteTag removes a tag by name. Returns 404 if the tag does not exist.
func (p *Postgres) DeleteTag(ctx context.Context, name string) error {
	cmd, err := p.pool.Exec(ctx, dbsql.DeleteTag, name)
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}
	return nil
}

// tagSortColumn validates and returns the sort column for tags.
// Returns empty string if sort is empty (no ordering).
func tagSortColumn(sort string) (string, error) {
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
