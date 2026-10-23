package store_test

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/FranciscoBarao/catalog/config"
	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/slug"
	"github.com/FranciscoBarao/catalog/internal/store"
	"github.com/FranciscoBarao/catalog/internal/store/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type BoardgameSuite struct {
	suite.Suite
	db    *postgres.Postgres
	admin *pgxpool.Pool
	cfg   *config.PostgresConfig
}

func (suite *BoardgameSuite) SetupSuite() {
	port, stop, err := startTestPostgres("catalog-store")
	suite.Require().NoError(err)
	suite.T().Cleanup(stop)

	_, filename, _, _ := runtime.Caller(0)
	migrationDir := filepath.Join(filepath.Dir(filename), "..", "..", "internal", "store", "postgres", "migrations")

	ctx := context.Background()
	suite.cfg = &config.PostgresConfig{
		Host:          "localhost",
		Username:      "postgres",
		Password:      "postgres",
		Port:          fmt.Sprintf("%d", port),
		Database:      "postgres",
		MigrationPath: migrationDir,
	}

	suite.db, err = postgres.Connect(ctx, suite.cfg)
	suite.Require().NoError(err)

	suite.admin, err = pgxpool.New(ctx, fmt.Sprintf(
		"host=localhost port=%s user=postgres password=postgres dbname=postgres sslmode=disable",
		suite.cfg.Port,
	))
	suite.Require().NoError(err)
}

func (suite *BoardgameSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Shutdown()
	}
	if suite.admin != nil {
		suite.admin.Close()
	}
}

func (suite *BoardgameSuite) TearDownTest() {
	_, err := suite.admin.Exec(context.Background(), `
		TRUNCATE TABLE
			boardgame_contributions,
			boardgame_mechanisms,
			boardgame_categories,
			boardgames,
			categories
		RESTART IDENTITY CASCADE`)
	suite.Require().NoError(err)
}

func (suite *BoardgameSuite) TestCreate_PlayersCheckConstraint() {
	_, err := suite.db.CreateBoardgame(context.Background(), boardgame.CreateBoardgameDTO{
		Name:       "Bad Players",
		Slug:       "bad-players",
		MinPlayers: 0,
		MaxPlayers: 4,
	})
	suite.Assert().ErrorContains(err, "database operation failed")
	suite.Assert().NotErrorIs(err, store.ErrAlreadyExists)
	suite.Assert().NotErrorIs(err, store.ErrNotFound)
}

func (suite *BoardgameSuite) TestCreate_DuplicatedSlug() {
	ctx := context.Background()
	_, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "duplicated",
		Slug:       slug.FromName("duplicated"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Require().NoError(err)

	_, err = suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "duplicated",
		Slug:       slug.FromName("duplicated"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Assert().ErrorIs(err, store.ErrAlreadyExists)
}

func (suite *BoardgameSuite) TestGet() {
	ctx := context.Background()
	created, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "name01",
		Slug:       slug.FromName("name01"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Require().NoError(err)

	bg, err := suite.db.GetBoardgameByID(ctx, created.ID)
	suite.Assert().NoError(err)
	suite.Assert().Equal(created.Name, bg.Name)
	suite.Assert().Equal(created.Slug, bg.Slug)
}

func (suite *BoardgameSuite) TestCreate_ScalarFieldsRoundTrip() {
	ctx := context.Background()
	input := boardgame.CreateBoardgameDTO{
		Name:          "Scalar Fields",
		Slug:          slug.FromName("Scalar Fields"),
		Description:   "A trading game",
		YearPublished: 1995,
		MinPlayers:    1,
		MaxPlayers:    4,
		MinPlayTime:   45,
		MaxPlayTime:   90,
	}

	created, err := suite.db.CreateBoardgame(ctx, input)
	suite.Require().NoError(err)
	suite.Assert().Equal("A trading game", created.Description)
	suite.Assert().Equal(1995, created.YearPublished)
	suite.Assert().Equal(45, created.MinPlayTime)
	suite.Assert().Equal(90, created.MaxPlayTime)
}

func (suite *BoardgameSuite) TestCreate_ScalarDefaultsWhenUnset() {
	ctx := context.Background()
	created, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "Unset Scalars",
		Slug:       slug.FromName("Unset Scalars"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Require().NoError(err)
	suite.Assert().Equal("", created.Description)
	suite.Assert().Equal(0, created.YearPublished)
	suite.Assert().Equal(0, created.MinPlayTime)
	suite.Assert().Equal(0, created.MaxPlayTime)
}

func (suite *BoardgameSuite) TestGetAll() {
	ctx := context.Background()
	_, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "name01",
		Slug:       slug.FromName("name01"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Require().NoError(err)
	_, err = suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "name02",
		Slug:       slug.FromName("name02"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Require().NoError(err)

	all, total, err := suite.db.GetAllBoardgames(ctx, listopt.NewQuery(listopt.WithSort(listopt.Sort{Column: "id", Order: "asc"})), false)
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 2)
	suite.Assert().Equal(2, total)
}

func (suite *BoardgameSuite) TestUpdate() {
	ctx := context.Background()
	created, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "updateName",
		Slug:       slug.FromName("updateName"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Require().NoError(err)

	readBg, err := suite.db.GetBoardgameByID(ctx, created.ID)
	suite.Assert().NoError(err)

	readBg.MaxPlayers = 6
	suite.Assert().NoError(suite.db.UpdateBoardgame(ctx, &readBg))

	updated, err := suite.db.GetBoardgameByID(ctx, created.ID)
	suite.Assert().NoError(err)
	suite.Assert().Equal(6, updated.MaxPlayers)
}

func (suite *BoardgameSuite) TestUpdateWithAssociations_ReplacesCategories() {
	ctx := context.Background()
	strategy := &category.Category{Slug: "strategy", Name: "Strategy"}
	economic := &category.Category{Slug: "economic", Name: "Economic"}
	suite.Require().NoError(suite.db.CreateCategory(ctx, strategy))
	suite.Require().NoError(suite.db.CreateCategory(ctx, economic))

	created, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "Catan",
		Slug:       "catan",
		MinPlayers: 2,
		MaxPlayers: 4,
		Categories: []uint{strategy.ID},
	})
	suite.Require().NoError(err)
	suite.Require().Len(created.Categories, 1)
	suite.Require().Equal("strategy", created.Categories[0].Slug)

	created.Description = "trading"
	next := []uint{economic.ID}
	suite.Require().NoError(suite.db.UpdateBoardgameWithAssociations(ctx, &created, boardgame.UpdateAssociations{
		Categories: &next,
	}))

	got, err := suite.db.GetBoardgameBySlug(ctx, "catan")
	suite.Assert().NoError(err)
	suite.Assert().Equal("trading", got.Description)
	suite.Require().Len(got.Categories, 1)
	suite.Assert().Equal("economic", got.Categories[0].Slug)
}

func (suite *BoardgameSuite) TestUpdateWithAssociations_NotFound() {
	next := []uint{}
	err := suite.db.UpdateBoardgameWithAssociations(context.Background(), &boardgame.Boardgame{
		ID:         999,
		Name:       "missing",
		Slug:       "missing",
		MinPlayers: 1,
		MaxPlayers: 2,
	}, boardgame.UpdateAssociations{Categories: &next})
	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *BoardgameSuite) TestUpdate_NotFound() {
	ctx := context.Background()
	err := suite.db.UpdateBoardgame(ctx, &boardgame.Boardgame{Name: "not found", Slug: "not-found", MinPlayers: 1, MaxPlayers: 2})
	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *BoardgameSuite) TestSoftDelete() {
	ctx := context.Background()
	created, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "deleteName",
		Slug:       slug.FromName("deleteName"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Require().NoError(err)

	suite.Assert().NoError(suite.db.DeleteBoardgame(ctx, created.ID, false))

	_, err = suite.db.GetBoardgameByID(ctx, created.ID)
	suite.Assert().Error(err)
}

func (suite *BoardgameSuite) TestDelete_NotFound() {
	ctx := context.Background()
	err := suite.db.DeleteBoardgame(ctx, 99999, false)
	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *BoardgameSuite) TestGetAll_FilterLike() {
	ctx := context.Background()
	_, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "Catan",
		Slug:       slug.FromName("Catan"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Require().NoError(err)
	_, err = suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "Vagrantsong",
		Slug:       slug.FromName("Vagrantsong"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Require().NoError(err)

	all, total, err := suite.db.GetAllBoardgames(ctx, listopt.NewQuery(listopt.WithFilter(listopt.Filter{Column: "name", Operator: listopt.Like, Value: "cat"})), false)
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 1)
	suite.Assert().Equal(1, total)
	suite.Assert().Equal("Catan", all[0].Name)
}

func (suite *BoardgameSuite) TestGetAll_FilterNumeric() {
	ctx := context.Background()
	input1 := boardgame.CreateBoardgameDTO{
		Name:       "Catan",
		Slug:       slug.FromName("Catan"),
		MinPlayers: 1,
		MaxPlayers: 4,
	}
	_, err := suite.db.CreateBoardgame(ctx, input1)
	suite.Require().NoError(err)

	input2 := boardgame.CreateBoardgameDTO{
		Name:       "Vagrantsong",
		Slug:       slug.FromName("Vagrantsong"),
		MinPlayers: 1,
		MaxPlayers: 2,
	}
	_, err = suite.db.CreateBoardgame(ctx, input2)
	suite.Require().NoError(err)

	all, total, err := suite.db.GetAllBoardgames(ctx, listopt.NewQuery(listopt.WithFilter(listopt.Filter{Column: "max_players", Operator: listopt.Lt, Value: "3", ValueKind: listopt.KindInt})), false)
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 1)
	suite.Assert().Equal(1, total)
	suite.Assert().Equal("Vagrantsong", all[0].Name)
}

func (suite *BoardgameSuite) TestGetAll_Pagination() {
	ctx := context.Background()
	for i := 0; i < 25; i++ {
		name := fmt.Sprintf("game-%02d", i)
		_, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
			Name:       name,
			Slug:       slug.FromName(name),
			MinPlayers: 1,
			MaxPlayers: 4,
		})
		suite.Require().NoError(err)
	}

	// Page 2 with pageSize 10 -> 10 items, total 25.
	page2, total, err := suite.db.GetAllBoardgames(ctx,
		listopt.NewQuery(listopt.WithSort(listopt.Sort{Column: "name", Order: "asc"}), listopt.WithPagination(listopt.Pagination{Page: 2, PageSize: 10})), false)
	suite.Assert().NoError(err)
	suite.Assert().Len(page2, 10)
	suite.Assert().Equal(25, total)
	suite.Assert().Equal("game-10", page2[0].Name)

	// Last page is partial (5 items).
	page3, total, err := suite.db.GetAllBoardgames(ctx,
		listopt.NewQuery(listopt.WithSort(listopt.Sort{Column: "name", Order: "asc"}), listopt.WithPagination(listopt.Pagination{Page: 3, PageSize: 10})), false)
	suite.Assert().NoError(err)
	suite.Assert().Len(page3, 5)
	suite.Assert().Equal(25, total)
}

func (suite *BoardgameSuite) TestGetAll_IncludeDeleted() {
	ctx := context.Background()
	_, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "Active",
		Slug:       slug.FromName("Active"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Require().NoError(err)
	deleted, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "Deleted",
		Slug:       slug.FromName("Deleted"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Require().NoError(err)
	suite.Require().NoError(suite.db.DeleteBoardgame(ctx, deleted.ID, false))

	active, total, err := suite.db.GetAllBoardgames(ctx, listopt.NewQuery(listopt.WithSort(listopt.Sort{Column: "name", Order: "asc"})), false)
	suite.Assert().NoError(err)
	suite.Assert().Len(active, 1)
	suite.Assert().Equal(1, total)
	suite.Assert().Equal("Active", active[0].Name)

	all, total, err := suite.db.GetAllBoardgames(ctx, listopt.NewQuery(listopt.WithSort(listopt.Sort{Column: "name", Order: "asc"})), true)
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 2)
	suite.Assert().Equal(2, total)
}

func (suite *BoardgameSuite) TestGetAll_MultipleFiltersAnd() {
	ctx := context.Background()
	catan := boardgame.CreateBoardgameDTO{
		Name:       "Catan",
		Slug:       slug.FromName("Catan"),
		MinPlayers: 1,
		MaxPlayers: 4,
	}
	_, err := suite.db.CreateBoardgame(ctx, catan)
	suite.Require().NoError(err)

	carcassonne := boardgame.CreateBoardgameDTO{
		Name:       "Carcassonne",
		Slug:       slug.FromName("Carcassonne"),
		MinPlayers: 1,
		MaxPlayers: 2,
	}
	_, err = suite.db.CreateBoardgame(ctx, carcassonne)
	suite.Require().NoError(err)

	vagrant := boardgame.CreateBoardgameDTO{
		Name:       "Vagrantsong",
		Slug:       slug.FromName("Vagrantsong"),
		MinPlayers: 1,
		MaxPlayers: 5,
	}
	_, err = suite.db.CreateBoardgame(ctx, vagrant)
	suite.Require().NoError(err)

	// name ILIKE '%ca%' AND max_players >= 3 -> only Catan.
	all, total, err := suite.db.GetAllBoardgames(ctx, listopt.NewQuery(
		listopt.WithFilter(listopt.Filter{Column: "name", Operator: listopt.Like, Value: "ca"}),
		listopt.WithFilter(listopt.Filter{Column: "max_players", Operator: listopt.Ge, Value: "3", ValueKind: listopt.KindInt}),
	), false)
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 1)
	suite.Assert().Equal(1, total)
	suite.Assert().Equal("Catan", all[0].Name)
}

func (suite *BoardgameSuite) TestGetAll_EqNumericColumn() {
	ctx := context.Background()
	four := boardgame.CreateBoardgameDTO{
		Name:       "Catan",
		Slug:       slug.FromName("Catan"),
		MinPlayers: 1,
		MaxPlayers: 4,
	}
	_, err := suite.db.CreateBoardgame(ctx, four)
	suite.Require().NoError(err)

	two := boardgame.CreateBoardgameDTO{
		Name:       "Carcassonne",
		Slug:       slug.FromName("Carcassonne"),
		MinPlayers: 1,
		MaxPlayers: 2,
	}
	_, err = suite.db.CreateBoardgame(ctx, two)
	suite.Require().NoError(err)

	all, total, err := suite.db.GetAllBoardgames(ctx,
		listopt.NewQuery(listopt.WithFilter(listopt.Filter{Column: "max_players", Operator: listopt.Eq, Value: "4", ValueKind: listopt.KindInt})), false)
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 1)
	suite.Assert().Equal(1, total)
	suite.Assert().Equal("Catan", all[0].Name)
}

func (suite *BoardgameSuite) TestGetAll_EqStringColumnNumericValue() {
	ctx := context.Background()
	_, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "123",
		Slug:       slug.FromName("123"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Require().NoError(err)
	_, err = suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "Catan",
		Slug:       slug.FromName("Catan"),
		MinPlayers: 1,
		MaxPlayers: 4,
	})
	suite.Require().NoError(err)

	// eq on the string "name" column with a numeric-looking value must be matched
	// as text, not coerced to an integer.
	all, total, err := suite.db.GetAllBoardgames(ctx,
		listopt.NewQuery(listopt.WithFilter(listopt.Filter{Column: "name", Operator: listopt.Eq, Value: "123"})), false)
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 1)
	suite.Assert().Equal(1, total)
	suite.Assert().Equal("123", all[0].Name)
}

func (suite *BoardgameSuite) TestGetAll_PageOutOfRange() {
	ctx := context.Background()
	for i := 0; i < 25; i++ {
		name := fmt.Sprintf("game-%02d", i)
		_, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
			Name:       name,
			Slug:       slug.FromName(name),
			MinPlayers: 1,
			MaxPlayers: 4,
		})
		suite.Require().NoError(err)
	}

	page, total, err := suite.db.GetAllBoardgames(ctx,
		listopt.NewQuery(listopt.WithSort(listopt.Sort{Column: "name", Order: "asc"}), listopt.WithPagination(listopt.Pagination{Page: 100, PageSize: 10})), false)
	suite.Assert().NoError(err)
	suite.Assert().Len(page, 0)
	suite.Assert().Equal(25, total)
}

func (suite *BoardgameSuite) TestConnect_NoPendingMigrations() {
	again, err := postgres.Connect(context.Background(), suite.cfg)
	suite.Require().NoError(err)
	again.Shutdown()
}

func TestBoardgameSuite(t *testing.T) {
	suite.Run(t, new(BoardgameSuite))
}
