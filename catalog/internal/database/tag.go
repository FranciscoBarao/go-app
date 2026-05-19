package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"

	dbsql "github.com/FranciscoBarao/catalog/internal/database/sql"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/middleware"
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
func (p *Postgres) GetAllTags(ctx context.Context, filter listopt.Params) ([]tag.Tag, error) {

	q := dbsql.SelectAllTags
	var args []any

	if where, arg := filterClause(filter); where != "" {
		q += where
		args = append(args, arg)
	}

	if filter.Sort.Column != "" {
		q += fmt.Sprintf(dbsql.OrderBy, filter.Sort.Column, filter.Sort.Order)
	}

	logging.FromCtx(ctx).Debug().Str("query", q).Msg("GetAllTags")

	rows, err := p.pool.Query(ctx, q, args...)
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
