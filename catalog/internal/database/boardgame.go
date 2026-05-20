package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/category"
	dbsql "github.com/FranciscoBarao/catalog/internal/database/sql"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/tag"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateBoardgame inserts a new boardgame and its associations within a transaction.
func (p *Postgres) CreateBoardgame(ctx context.Context, bg *boardgame.Boardgame) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return mapPgError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = tx.QueryRow(ctx, dbsql.InsertBoardgame,
		bg.Name, bg.Publisher, bg.PlayerNumber, bg.BoardgameID,
	).Scan(&bg.ID)
	if err != nil {
		return mapPgError(err)
	}

	for _, t := range bg.Tags {
		if _, err := tx.Exec(ctx, dbsql.InsertBoardgameTag, bg.ID, t.Name); err != nil {
			return mapPgError(err)
		}
	}

	for _, c := range bg.Categories {
		if _, err := tx.Exec(ctx, dbsql.InsertBoardgameCategory, bg.ID, c.Name); err != nil {
			return mapPgError(err)
		}
	}

	for _, m := range bg.Mechanisms {
		if _, err := tx.Exec(ctx, dbsql.InsertBoardgameMechanism, bg.ID, m.Name); err != nil {
			return mapPgError(err)
		}
	}

	return tx.Commit(ctx)
}

// GetBoardgameByID retrieves a single boardgame by ID with all associations loaded.
func (p *Postgres) GetBoardgameByID(ctx context.Context, id uint) (boardgame.Boardgame, error) {
	var bg boardgame.Boardgame
	err := p.pool.QueryRow(ctx, dbsql.SelectBoardgameByID, id).Scan(
		&bg.ID, &bg.CreatedAt, &bg.UpdatedAt,
		&bg.Name, &bg.Publisher, &bg.PlayerNumber, &bg.BoardgameID,
	)
	if err != nil {
		return bg, mapPgError(err)
	}

	if err := loadBoardgameAssociations(ctx, p.pool, &bg); err != nil {
		return bg, err
	}

	return bg, nil
}

// GetAllBoardgames retrieves all boardgames with optional filtering and sorting.
func (p *Postgres) GetAllBoardgames(ctx context.Context, filter listopt.Params) ([]boardgame.Boardgame, error) {

	q := dbsql.SelectAllBoardgames
	var args []any

	if where, arg := filterClause(filter); where != "" {
		q += where
		args = append(args, arg)
	}

	if filter.Sort.Column != "" {
		q += fmt.Sprintf(dbsql.OrderBy, filter.Sort.Column, filter.Sort.Order)
	}

	logging.FromCtx(ctx).Debug().Str("query", q).Msg("GetAllBoardgames")

	rows, err := p.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, mapPgError(err)
	}

	boardgames, err := pgx.CollectRows(rows, scanBoardgameRow)
	if err != nil {
		return nil, mapPgError(err)
	}

	for i := range boardgames {
		if err := loadBoardgameAssociations(ctx, p.pool, &boardgames[i]); err != nil {
			return nil, err
		}
	}

	return boardgames, nil
}

// UpdateBoardgame updates a boardgame's mutable fields.
func (p *Postgres) UpdateBoardgame(ctx context.Context, bg *boardgame.Boardgame) error {
	cmd, err := p.pool.Exec(ctx, dbsql.UpdateBoardgame,
		bg.Name, bg.Publisher, bg.PlayerNumber, bg.BoardgameID, bg.ID,
	)
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}
	return nil
}

// UpdateBoardgameWithAssociations updates a boardgame and conditionally replaces associations in a single transaction.
func (p *Postgres) UpdateBoardgameWithAssociations(ctx context.Context, bg *boardgame.Boardgame, assoc boardgame.UpdateAssociations) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return mapPgError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	cmd, err := tx.Exec(ctx, dbsql.UpdateBoardgame, bg.Name, bg.Publisher, bg.PlayerNumber, bg.BoardgameID, bg.ID)
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}

	if assoc.Tags != nil {
		if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameTags, bg.ID); err != nil {
			return mapPgError(err)
		}
		for _, t := range *assoc.Tags {
			if _, err := tx.Exec(ctx, dbsql.InsertBoardgameTag, bg.ID, t.Name); err != nil {
				return mapPgError(err)
			}
		}
	}

	if assoc.Categories != nil {
		if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameCategories, bg.ID); err != nil {
			return mapPgError(err)
		}
		for _, c := range *assoc.Categories {
			if _, err := tx.Exec(ctx, dbsql.InsertBoardgameCategory, bg.ID, c.Name); err != nil {
				return mapPgError(err)
			}
		}
	}

	if assoc.Mechanisms != nil {
		if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameMechanisms, bg.ID); err != nil {
			return mapPgError(err)
		}
		for _, m := range *assoc.Mechanisms {
			if _, err := tx.Exec(ctx, dbsql.InsertBoardgameMechanism, bg.ID, m.Name); err != nil {
				return mapPgError(err)
			}
		}
	}

	return tx.Commit(ctx)
}

// DeleteBoardgame removes a boardgame and its join table entries within a transaction.
func (p *Postgres) DeleteBoardgame(ctx context.Context, id uint) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return mapPgError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameTags, id); err != nil {
		return mapPgError(err)
	}
	if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameCategories, id); err != nil {
		return mapPgError(err)
	}
	if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameMechanisms, id); err != nil {
		return mapPgError(err)
	}

	cmd, err := tx.Exec(ctx, dbsql.DeleteBoardgame, id)
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}

	return tx.Commit(ctx)
}

// loadBoardgameAssociations loads all associations for a boardgame.
func loadBoardgameAssociations(ctx context.Context, pool *pgxpool.Pool, bg *boardgame.Boardgame) error {
	var err error

	bg.Tags, err = loadBoardgameTags(ctx, pool, bg.ID)
	if err != nil {
		return err
	}
	bg.Categories, err = loadBoardgameCategories(ctx, pool, bg.ID)
	if err != nil {
		return err
	}
	bg.Mechanisms, err = loadBoardgameMechanisms(ctx, pool, bg.ID)
	if err != nil {
		return err
	}

	bg.Expansions, err = loadBoardgameExpansions(ctx, pool, bg.ID)
	return err
}

func loadBoardgameTags(ctx context.Context, pool *pgxpool.Pool, boardgameID uint) ([]tag.Tag, error) {
	rows, err := pool.Query(ctx, dbsql.SelectBoardgameTags, boardgameID)
	if err != nil {
		return nil, mapPgError(err)
	}
	tags, err := pgx.CollectRows(rows, pgx.RowToStructByPos[tag.Tag])
	return tags, mapPgError(err)
}

func loadBoardgameCategories(ctx context.Context, pool *pgxpool.Pool, boardgameID uint) ([]category.Category, error) {
	rows, err := pool.Query(ctx, dbsql.SelectBoardgameCategories, boardgameID)
	if err != nil {
		return nil, mapPgError(err)
	}
	categories, err := pgx.CollectRows(rows, pgx.RowToStructByPos[category.Category])
	return categories, mapPgError(err)
}

func loadBoardgameMechanisms(ctx context.Context, pool *pgxpool.Pool, boardgameID uint) ([]mechanism.Mechanism, error) {
	rows, err := pool.Query(ctx, dbsql.SelectBoardgameMechanisms, boardgameID)
	if err != nil {
		return nil, mapPgError(err)
	}
	mechanisms, err := pgx.CollectRows(rows, pgx.RowToStructByPos[mechanism.Mechanism])
	return mechanisms, mapPgError(err)
}

func loadBoardgameExpansions(ctx context.Context, pool *pgxpool.Pool, boardgameID uint) ([]boardgame.Boardgame, error) {
	rows, err := pool.Query(ctx, dbsql.SelectBoardgameExpansions, boardgameID)
	if err != nil {
		return nil, mapPgError(err)
	}
	expansions, err := pgx.CollectRows(rows, scanBoardgameRow)
	return expansions, mapPgError(err)
}

// scanBoardgameRow scans a boardgame row into a Boardgame struct (without associations).
func scanBoardgameRow(row pgx.CollectableRow) (boardgame.Boardgame, error) {
	var bg boardgame.Boardgame
	err := row.Scan(
		&bg.ID, &bg.CreatedAt, &bg.UpdatedAt,
		&bg.Name, &bg.Publisher, &bg.PlayerNumber, &bg.BoardgameID,
	)
	return bg, err
}
