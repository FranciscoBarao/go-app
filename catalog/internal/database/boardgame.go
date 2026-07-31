package database

import (
	"context"
	"net/http"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/contributor"
	dbsql "github.com/FranciscoBarao/catalog/internal/database/sql"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateBoardgame inserts a new boardgame and its associations within a transaction.
func (p *Postgres) CreateBoardgame(ctx context.Context, input boardgame.CreateBoardgameDTO) (boardgame.Boardgame, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return boardgame.Boardgame{}, mapPgError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var id uint
	err = tx.QueryRow(ctx, dbsql.InsertBoardgame,
		input.Slug, input.Name, input.Description, input.YearPublished,
		input.MinPlayers, input.MaxPlayers, input.MinPlayTime, input.MaxPlayTime,
		input.MinAge, input.BggID, input.ParentID,
	).Scan(&id)
	if err != nil {
		return boardgame.Boardgame{}, mapPgError(err)
	}

	if err := insertBoardgameCategories(ctx, tx, id, input.Categories); err != nil {
		return boardgame.Boardgame{}, err
	}
	if err := insertBoardgameMechanisms(ctx, tx, id, input.Mechanisms); err != nil {
		return boardgame.Boardgame{}, err
	}
	if err := insertBoardgameContributions(ctx, tx, id, input.Contributions); err != nil {
		return boardgame.Boardgame{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return boardgame.Boardgame{}, mapPgError(err)
	}

	return p.getBoardgame(ctx, dbsql.SelectBoardgameByID, id)
}

// GetBoardgameByID retrieves a single active boardgame by ID.
func (p *Postgres) GetBoardgameByID(ctx context.Context, id uint) (boardgame.Boardgame, error) {
	return p.getBoardgame(ctx, dbsql.SelectBoardgameByID, id)
}

// GetBoardgameBySlug retrieves a single active boardgame by slug.
func (p *Postgres) GetBoardgameBySlug(ctx context.Context, slug string) (boardgame.Boardgame, error) {
	return p.getBoardgame(ctx, dbsql.SelectBoardgameBySlug, slug)
}

func (p *Postgres) getBoardgame(ctx context.Context, query string, arg any) (boardgame.Boardgame, error) {
	var bg boardgame.Boardgame
	err := p.pool.QueryRow(ctx, query, arg).Scan(
		&bg.ID, &bg.Slug, &bg.CreatedAt, &bg.UpdatedAt, &bg.DeletedAt,
		&bg.Name, &bg.Description, &bg.YearPublished,
		&bg.MinPlayers, &bg.MaxPlayers, &bg.MinPlayTime, &bg.MaxPlayTime,
		&bg.MinAge, &bg.BggID, &bg.BoardgameID,
	)
	if err != nil {
		return bg, mapPgError(err)
	}

	if err := loadBoardgameAssociations(ctx, p.pool, &bg); err != nil {
		return bg, err
	}

	return bg, nil
}

// GetAllBoardgames retrieves active boardgames with optional filtering, sorting,
// and pagination. It returns the page of results and the total count of matching
// rows (before pagination).
func (p *Postgres) GetAllBoardgames(ctx context.Context, filter listopt.Params, includeDeleted bool) ([]boardgame.Boardgame, int, error) {
	selectBase := dbsql.SelectAllBoardgames
	countBase := dbsql.CountBoardgames
	if includeDeleted {
		selectBase = dbsql.SelectAllBoardgamesIncludingDeleted
		countBase = dbsql.CountBoardgamesIncludingDeleted
	}

	countQuery, countArgs := buildCountQuery(countBase, filter)
	selectQuery, selectArgs := buildPaginatedQuery(selectBase, filter)

	logging.FromCtx(ctx).Debug().Str("query", selectQuery).Msg("GetAllBoardgames")

	var total int
	if err := p.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, mapPgError(err)
	}

	rows, err := p.pool.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, mapPgError(err)
	}

	boardgames, err := pgx.CollectRows(rows, scanBoardgameRow)
	if err != nil {
		return nil, 0, mapPgError(err)
	}

	for i := range boardgames {
		if err := loadBoardgameAssociations(ctx, p.pool, &boardgames[i]); err != nil {
			return nil, 0, err
		}
	}

	return boardgames, total, nil
}

// UpdateBoardgame updates a boardgame's mutable fields.
func (p *Postgres) UpdateBoardgame(ctx context.Context, bg *boardgame.Boardgame) error {
	cmd, err := p.pool.Exec(ctx, dbsql.UpdateBoardgame,
		bg.Name, bg.Description, bg.YearPublished,
		bg.MinPlayers, bg.MaxPlayers, bg.MinPlayTime, bg.MaxPlayTime,
		bg.MinAge, bg.BggID, bg.BoardgameID, bg.ID,
	)
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}
	return nil
}

// UpdateBoardgameWithAssociations updates a boardgame and conditionally replaces associations.
func (p *Postgres) UpdateBoardgameWithAssociations(ctx context.Context, bg *boardgame.Boardgame, assoc boardgame.UpdateAssociations) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return mapPgError(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	cmd, err := tx.Exec(ctx, dbsql.UpdateBoardgame,
		bg.Name, bg.Description, bg.YearPublished,
		bg.MinPlayers, bg.MaxPlayers, bg.MinPlayTime, bg.MaxPlayTime,
		bg.MinAge, bg.BggID, bg.BoardgameID, bg.ID,
	)
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}

	if assoc.Categories != nil {
		if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameCategories, bg.ID); err != nil {
			return mapPgError(err)
		}
		if err := insertBoardgameCategories(ctx, tx, bg.ID, *assoc.Categories); err != nil {
			return err
		}
	}

	if assoc.Mechanisms != nil {
		if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameMechanisms, bg.ID); err != nil {
			return mapPgError(err)
		}
		if err := insertBoardgameMechanisms(ctx, tx, bg.ID, *assoc.Mechanisms); err != nil {
			return err
		}
	}

	if assoc.Contributions != nil {
		if _, err := tx.Exec(ctx, dbsql.DeleteBoardgameContributions, bg.ID); err != nil {
			return mapPgError(err)
		}
		if err := insertBoardgameContributions(ctx, tx, bg.ID, *assoc.Contributions); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// DeleteBoardgame soft-deletes or hard-deletes a boardgame.
func (p *Postgres) DeleteBoardgame(ctx context.Context, id uint, hard bool) error {
	var (
		cmd pgconn.CommandTag
		err error
	)
	if hard {
		tx, txErr := p.pool.Begin(ctx)
		if txErr != nil {
			return mapPgError(txErr)
		}
		defer func() { _ = tx.Rollback(ctx) }()

		if _, err = tx.Exec(ctx, dbsql.DeleteBoardgameCategories, id); err != nil {
			return mapPgError(err)
		}
		if _, err = tx.Exec(ctx, dbsql.DeleteBoardgameMechanisms, id); err != nil {
			return mapPgError(err)
		}
		if _, err = tx.Exec(ctx, dbsql.DeleteBoardgameContributions, id); err != nil {
			return mapPgError(err)
		}
		if _, err = tx.Exec(ctx, dbsql.HardDeleteBoardgameExpansions, id); err != nil {
			return mapPgError(err)
		}
		cmd, err = tx.Exec(ctx, dbsql.HardDeleteBoardgame, id)
		if err != nil {
			return mapPgError(err)
		}
		if cmd.RowsAffected() == 0 {
			return middleware.NewError(http.StatusNotFound, "record not found")
		}
		return tx.Commit(ctx)
	}

	tx, txErr := p.pool.Begin(ctx)
	if txErr != nil {
		return mapPgError(txErr)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err = tx.Exec(ctx, dbsql.SoftDeleteBoardgameExpansions, id); err != nil {
		return mapPgError(err)
	}
	cmd, err = tx.Exec(ctx, dbsql.SoftDeleteBoardgame, id)
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}
	return mapPgError(tx.Commit(ctx))
}

func insertBoardgameCategories(ctx context.Context, tx pgx.Tx, boardgameID uint, categoryIDs []uint) error {
	for _, categoryID := range categoryIDs {
		if _, err := tx.Exec(ctx, dbsql.InsertBoardgameCategory, boardgameID, categoryID); err != nil {
			return mapPgError(err)
		}
	}
	return nil
}

func insertBoardgameMechanisms(ctx context.Context, tx pgx.Tx, boardgameID uint, mechanismIDs []uint) error {
	for _, mechanismID := range mechanismIDs {
		if _, err := tx.Exec(ctx, dbsql.InsertBoardgameMechanism, boardgameID, mechanismID); err != nil {
			return mapPgError(err)
		}
	}
	return nil
}

func insertBoardgameContributions(ctx context.Context, tx pgx.Tx, boardgameID uint, contributions []boardgame.ContributionDTO) error {
	for _, contrib := range contributions {
		if _, err := tx.Exec(ctx, dbsql.InsertBoardgameContribution,
			boardgameID, contrib.ContributorID, contrib.Role, contrib.CreditOrder,
		); err != nil {
			return mapPgError(err)
		}
	}
	return nil
}

func loadBoardgameAssociations(ctx context.Context, pool *pgxpool.Pool, bg *boardgame.Boardgame) error {
	var err error

	bg.Categories, err = loadBoardgameCategories(ctx, pool, bg.ID)
	if err != nil {
		return err
	}
	bg.Mechanisms, err = loadBoardgameMechanisms(ctx, pool, bg.ID)
	if err != nil {
		return err
	}
	bg.Contributions, err = loadBoardgameContributions(ctx, pool, bg.ID)
	if err != nil {
		return err
	}

	bg.Expansions, err = loadBoardgameExpansions(ctx, pool, bg.ID)
	return err
}

func loadBoardgameCategories(ctx context.Context, pool *pgxpool.Pool, boardgameID uint) ([]category.Category, error) {
	rows, err := pool.Query(ctx, dbsql.SelectBoardgameCategories, boardgameID)
	if err != nil {
		return nil, mapPgError(err)
	}
	categories, err := pgx.CollectRows(rows, scanCategoryRow)
	return categories, mapPgError(err)
}

func loadBoardgameMechanisms(ctx context.Context, pool *pgxpool.Pool, boardgameID uint) ([]mechanism.Mechanism, error) {
	rows, err := pool.Query(ctx, dbsql.SelectBoardgameMechanisms, boardgameID)
	if err != nil {
		return nil, mapPgError(err)
	}
	mechanisms, err := pgx.CollectRows(rows, scanMechanismRow)
	return mechanisms, mapPgError(err)
}

func loadBoardgameContributions(ctx context.Context, pool *pgxpool.Pool, boardgameID uint) ([]contributor.Contribution, error) {
	rows, err := pool.Query(ctx, dbsql.SelectBoardgameContributions, boardgameID)
	if err != nil {
		return nil, mapPgError(err)
	}
	defer rows.Close()

	var contributions []contributor.Contribution
	for rows.Next() {
		var c contributor.Contributor
		var role string
		var creditOrder *int
		if err := rows.Scan(
			&c.ID, &c.Slug, &c.Name, &c.Bio, &c.BggID, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
			&role, &creditOrder,
		); err != nil {
			return nil, mapPgError(err)
		}
		contributions = append(contributions, contributor.Contribution{
			Contributor: c,
			Role:        contributor.Role(role),
			CreditOrder: creditOrder,
		})
	}
	return contributions, mapPgError(rows.Err())
}

func loadBoardgameExpansions(ctx context.Context, pool *pgxpool.Pool, boardgameID uint) ([]boardgame.Boardgame, error) {
	rows, err := pool.Query(ctx, dbsql.SelectBoardgameExpansions, boardgameID)
	if err != nil {
		return nil, mapPgError(err)
	}
	expansions, err := pgx.CollectRows(rows, scanBoardgameRow)
	return expansions, mapPgError(err)
}

func scanBoardgameRow(row pgx.CollectableRow) (boardgame.Boardgame, error) {
	var bg boardgame.Boardgame
	err := row.Scan(
		&bg.ID, &bg.Slug, &bg.CreatedAt, &bg.UpdatedAt, &bg.DeletedAt,
		&bg.Name, &bg.Description, &bg.YearPublished,
		&bg.MinPlayers, &bg.MaxPlayers, &bg.MinPlayTime, &bg.MaxPlayTime,
		&bg.MinAge, &bg.BggID, &bg.BoardgameID,
	)
	return bg, err
}
