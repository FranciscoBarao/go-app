package store_test

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/FranciscoBarao/catalog/config"
	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/store"
	"github.com/FranciscoBarao/catalog/internal/store/postgres"
)

type CategorySuite struct {
	suite.Suite
	db    *postgres.Postgres
	admin *pgxpool.Pool
}

func (suite *CategorySuite) SetupSuite() {
	port, stop, err := startTestPostgres("catalog-category")
	suite.Require().NoError(err)
	suite.T().Cleanup(stop)

	_, filename, _, _ := runtime.Caller(0)
	migrationDir := filepath.Join(filepath.Dir(filename), "..", "..", "internal", "store", "postgres", "migrations")

	ctx := context.Background()
	cfg := &config.PostgresConfig{
		Host:          "localhost",
		Username:      "postgres",
		Password:      "postgres",
		Port:          fmt.Sprintf("%d", port),
		Database:      "postgres",
		MigrationPath: migrationDir,
	}

	suite.db, err = postgres.Connect(ctx, cfg)
	suite.Require().NoError(err)

	suite.admin, err = pgxpool.New(ctx, fmt.Sprintf(
		"host=localhost port=%s user=postgres password=postgres dbname=postgres sslmode=disable",
		cfg.Port,
	))
	suite.Require().NoError(err)
}

func (suite *CategorySuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Shutdown()
	}
	if suite.admin != nil {
		suite.admin.Close()
	}
}

func (suite *CategorySuite) TearDownTest() {
	_, err := suite.admin.Exec(context.Background(), `
		TRUNCATE TABLE
			boardgame_categories,
			boardgames,
			categories
		RESTART IDENTITY CASCADE`)
	suite.Require().NoError(err)
}

func (suite *CategorySuite) TestRoundTripAndList() {
	ctx := context.Background()
	bggID := 1020
	created := &category.Category{Slug: "strategy", Name: "Strategy", BggID: &bggID}
	suite.Require().NoError(suite.db.CreateCategory(ctx, created))
	suite.Require().NotZero(created.ID)

	got, err := suite.db.GetCategoryBySlug(ctx, "strategy")
	suite.Assert().NoError(err)
	suite.Assert().Equal(created.ID, got.ID)
	suite.Assert().Equal("Strategy", got.Name)
	suite.Require().NotNil(got.BggID)
	suite.Assert().Equal(bggID, *got.BggID)

	id, err := suite.db.GetCategoryIDBySlug(ctx, "strategy")
	suite.Assert().NoError(err)
	suite.Assert().Equal(created.ID, id)

	err = suite.db.CreateCategory(ctx, &category.Category{Slug: "economic", Name: "Economic"})
	suite.Require().NoError(err)

	listed, total, err := suite.db.GetAllCategories(ctx, listopt.NewQuery(listopt.WithFilter(listopt.Filter{
		Column:   "name",
		Operator: listopt.Like,
		Value:    "strat",
	})))
	suite.Assert().NoError(err)
	suite.Assert().Equal(1, total)
	suite.Require().Len(listed, 1)
	suite.Assert().Equal("strategy", listed[0].Slug)
}

func (suite *CategorySuite) TestDuplicateSlug() {
	ctx := context.Background()
	suite.Require().NoError(suite.db.CreateCategory(ctx, &category.Category{Slug: "strategy", Name: "Strategy"}))

	err := suite.db.CreateCategory(ctx, &category.Category{Slug: "strategy", Name: "Strategy Again"})
	suite.Assert().ErrorIs(err, store.ErrAlreadyExists)
}

func (suite *CategorySuite) TestDuplicateBggID() {
	ctx := context.Background()
	bggID := 55
	suite.Require().NoError(suite.db.CreateCategory(ctx, &category.Category{Slug: "strategy", Name: "Strategy", BggID: &bggID}))

	err := suite.db.CreateCategory(ctx, &category.Category{Slug: "economic", Name: "Economic", BggID: &bggID})
	suite.Assert().ErrorIs(err, store.ErrAlreadyExists)
}

func (suite *CategorySuite) TestEmptySlug() {
	err := suite.db.CreateCategory(context.Background(), &category.Category{Slug: "", Name: "Blank"})
	suite.Assert().ErrorContains(err, "database operation failed")
	suite.Assert().NotErrorIs(err, store.ErrAlreadyExists)
}

func (suite *CategorySuite) TestGetMissing() {
	_, err := suite.db.GetCategoryBySlug(context.Background(), "missing")
	suite.Assert().ErrorIs(err, store.ErrNotFound)

	_, err = suite.db.GetCategoryIDBySlug(context.Background(), "missing")
	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *CategorySuite) TestSoftDeleteHidesRowAndKeepsSlug() {
	ctx := context.Background()
	suite.Require().NoError(suite.db.CreateCategory(ctx, &category.Category{Slug: "strategy", Name: "Strategy"}))
	suite.Require().NoError(suite.db.DeleteCategory(ctx, "strategy", false))

	_, err := suite.db.GetCategoryBySlug(ctx, "strategy")
	suite.Assert().ErrorIs(err, store.ErrNotFound)

	err = suite.db.CreateCategory(ctx, &category.Category{Slug: "strategy", Name: "Strategy"})
	suite.Assert().ErrorIs(err, store.ErrAlreadyExists)
}

func (suite *CategorySuite) TestHardDeleteFreesSlug() {
	ctx := context.Background()
	suite.Require().NoError(suite.db.CreateCategory(ctx, &category.Category{Slug: "strategy", Name: "Strategy"}))
	suite.Require().NoError(suite.db.DeleteCategory(ctx, "strategy", true))

	err := suite.db.CreateCategory(ctx, &category.Category{Slug: "strategy", Name: "Strategy"})
	suite.Assert().NoError(err)
}

func (suite *CategorySuite) TestDeleteMissing() {
	err := suite.db.DeleteCategory(context.Background(), "missing", false)
	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *CategorySuite) TestHardDeleteCascadesJoin() {
	ctx := context.Background()
	cat := &category.Category{Slug: "strategy", Name: "Strategy"}
	suite.Require().NoError(suite.db.CreateCategory(ctx, cat))

	created, err := suite.db.CreateBoardgame(ctx, boardgame.CreateBoardgameDTO{
		Name:       "Catan",
		Slug:       "catan",
		MinPlayers: 2,
		MaxPlayers: 4,
		Categories: []uint{cat.ID},
	})
	suite.Require().NoError(err)
	suite.Require().Len(created.Categories, 1)

	var links int
	suite.Require().NoError(suite.admin.QueryRow(ctx, `SELECT COUNT(*) FROM boardgame_categories WHERE category_id = $1`, cat.ID).Scan(&links))
	suite.Require().Equal(1, links)

	suite.Require().NoError(suite.db.DeleteCategory(ctx, "strategy", true))

	suite.Require().NoError(suite.admin.QueryRow(ctx, `SELECT COUNT(*) FROM boardgame_categories WHERE category_id = $1`, cat.ID).Scan(&links))
	suite.Assert().Equal(0, links)

	got, err := suite.db.GetBoardgameBySlug(ctx, "catan")
	suite.Assert().NoError(err)
	suite.Assert().Empty(got.Categories)
}

func TestCategorySuite(t *testing.T) {
	suite.Run(t, new(CategorySuite))
}
