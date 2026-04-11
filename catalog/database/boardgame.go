package database

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FranciscoBarao/catalog/boardgame"
	"github.com/FranciscoBarao/catalog/category"
	dbsql "github.com/FranciscoBarao/catalog/database/sql"
	"github.com/FranciscoBarao/catalog/mechanism"
	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/tag"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateBoardgame inserts a new boardgame and its associations within a transaction.
func (p *Postgres) CreateBoardgame(ctx context.Context, bg *boardgame.Boardgame) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return mapPgError(err)
	}
	defer tx.Rollback(ctx)

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
func (p *Postgres) GetAllBoardgames(ctx context.Context, sort string) ([]boardgame.Boardgame, error) {
	col, err := boardgameSortColumn(sort)
	if err != nil {
		return nil, err
	}

	var rows pgx.Rows

	query := dbsql.SelectAllBoardgames
	if col != "" {
		query += fmt.Sprintf(dbsql.OrderBy, col)
	}

	rows, err = p.pool.Query(ctx, query)
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

// DeleteBoardgame removes a boardgame and its join table entries within a transaction.
func (p *Postgres) DeleteBoardgame(ctx context.Context, id uint) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return mapPgError(err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameTags, id); err != nil {
		return mapPgError(err)
	}
	if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameCategories, id); err != nil {
		return mapPgError(err)
	}
	if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameMechanisms, id); err != nil {
		return mapPgError(err)
	}
	if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameRatings, id); err != nil {
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

// ReplaceBoardgameTags replaces all tag associations for a boardgame within a transaction.
func (p *Postgres) ReplaceBoardgameTags(ctx context.Context, boardgameID uint, tags []tag.Tag) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return mapPgError(err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameTags, boardgameID); err != nil {
		return mapPgError(err)
	}

	for _, t := range tags {
		if _, err := tx.Exec(ctx, dbsql.InsertBoardgameTag, boardgameID, t.Name); err != nil {
			return mapPgError(err)
		}
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
	bg.Ratings, err = loadBoardgameRatings(ctx, pool, bg.ID)
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

func loadBoardgameRatings(ctx context.Context, pool *pgxpool.Pool, boardgameID uint) ([]boardgame.Rating, error) {
	rows, err := pool.Query(ctx, dbsql.SelectBoardgameRatings, boardgameID)
	if err != nil {
		return nil, mapPgError(err)
	}
	ratings, err := pgx.CollectRows(rows, pgx.RowToStructByPos[boardgame.Rating])
	return ratings, mapPgError(err)
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

// boardgameSortColumn validates and returns the sort column for boardgames.
func boardgameSortColumn(sort string) (string, error) {
	if sort == "" {
		return "", nil
	}
	switch sort {
	case "id", "name", "publisher", "player_number", "created_at", "updated_at":
		return sort, nil
	default:
		return "", middleware.NewError(http.StatusBadRequest, fmt.Sprintf("invalid sort column: %s", sort))
	}
}
