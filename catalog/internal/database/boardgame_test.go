package database

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/FranciscoBarao/catalog/config"
	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/slug"
	"github.com/stretchr/testify/suite"
)

type BoardgameSuite struct {
	suite.Suite
	db       *TestPostgres
	postgres *Postgres
}

func migrationDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "migrations")
}

func testCreateInput(name string) boardgame.CreateBoardgameDTO {
	return boardgame.CreateBoardgameDTO{
		Name:       name,
		Slug:       slug.FromName(name),
		MinPlayers: 1,
		MaxPlayers: 4,
	}
}

func (suite *BoardgameSuite) SetupSuite() {
	var err error
	suite.db, err = NewTestPostgres("test")
	suite.Require().NoError(err)

	ctx := context.Background()
	cfg := &config.PostgresConfig{
		Host:          "localhost",
		Username:      "postgres",
		Password:      "postgres",
		Port:          fmt.Sprintf("%d", suite.db.Port),
		Database:      "postgres",
		MigrationPath: migrationDir(),
	}

	suite.postgres, err = Connect(ctx, cfg)
	suite.Require().NoError(err)
}

func (suite *BoardgameSuite) TearDownSuite() {
	if suite.postgres != nil {
		suite.postgres.Close()
	}
	suite.db.Shutdown()
}

func (suite *BoardgameSuite) TearDownTest() {
	_, err := suite.postgres.pool.Exec(context.Background(), `TRUNCATE TABLE boardgames CASCADE;`)
	suite.Require().NoError(err)
}

func (suite *BoardgameSuite) TestCreate_DuplicatedSlug() {
	ctx := context.Background()
	_, err := suite.postgres.CreateBoardgame(ctx, testCreateInput("duplicated"))
	suite.Require().NoError(err)

	_, err = suite.postgres.CreateBoardgame(ctx, testCreateInput("duplicated"))
	suite.Assert().Error(err)
}

func (suite *BoardgameSuite) TestGet() {
	ctx := context.Background()
	created, err := suite.postgres.CreateBoardgame(ctx, testCreateInput("name01"))
	suite.Require().NoError(err)

	bg, err := suite.postgres.GetBoardgameByID(ctx, created.ID)
	suite.Assert().NoError(err)
	suite.Assert().Equal(created.Name, bg.Name)
	suite.Assert().Equal(created.Slug, bg.Slug)
}

func (suite *BoardgameSuite) TestCreate_ScalarFieldsRoundTrip() {
	ctx := context.Background()
	input := testCreateInput("Scalar Fields")
	input.Description = "A trading game"
	input.YearPublished = 1995
	input.MinPlayTime = 45
	input.MaxPlayTime = 90

	created, err := suite.postgres.CreateBoardgame(ctx, input)
	suite.Require().NoError(err)
	suite.Assert().Equal("A trading game", created.Description)
	suite.Assert().Equal(1995, created.YearPublished)
	suite.Assert().Equal(45, created.MinPlayTime)
	suite.Assert().Equal(90, created.MaxPlayTime)
}

func (suite *BoardgameSuite) TestCreate_ScalarDefaultsWhenUnset() {
	ctx := context.Background()
	created, err := suite.postgres.CreateBoardgame(ctx, testCreateInput("Unset Scalars"))
	suite.Require().NoError(err)
	suite.Assert().Equal("", created.Description)
	suite.Assert().Equal(0, created.YearPublished)
	suite.Assert().Equal(0, created.MinPlayTime)
	suite.Assert().Equal(0, created.MaxPlayTime)
}

func (suite *BoardgameSuite) TestGetAll() {
	ctx := context.Background()
	_, err := suite.postgres.CreateBoardgame(ctx, testCreateInput("name01"))
	suite.Require().NoError(err)
	_, err = suite.postgres.CreateBoardgame(ctx, testCreateInput("name02"))
	suite.Require().NoError(err)

	all, err := suite.postgres.GetAllBoardgames(ctx, listopt.Params{Sort: listopt.Sort{Column: "id", Order: "asc"}}, false)
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 2)
}

func (suite *BoardgameSuite) TestUpdate() {
	ctx := context.Background()
	created, err := suite.postgres.CreateBoardgame(ctx, testCreateInput("updateName"))
	suite.Require().NoError(err)

	readBg, err := suite.postgres.GetBoardgameByID(ctx, created.ID)
	suite.Assert().NoError(err)

	readBg.MaxPlayers = 6
	suite.Assert().NoError(suite.postgres.UpdateBoardgame(ctx, &readBg))

	updated, err := suite.postgres.GetBoardgameByID(ctx, created.ID)
	suite.Assert().NoError(err)
	suite.Assert().Equal(6, updated.MaxPlayers)
}

func (suite *BoardgameSuite) TestUpdate_NotFound() {
	ctx := context.Background()
	err := suite.postgres.UpdateBoardgame(ctx, &boardgame.Boardgame{Name: "not found", Slug: "not-found", MinPlayers: 1, MaxPlayers: 2})
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().True(errors.As(err, &mr))
	suite.Assert().Equal(http.StatusNotFound, mr.GetStatus())
}

func (suite *BoardgameSuite) TestSoftDelete() {
	ctx := context.Background()
	created, err := suite.postgres.CreateBoardgame(ctx, testCreateInput("deleteName"))
	suite.Require().NoError(err)

	suite.Assert().NoError(suite.postgres.DeleteBoardgame(ctx, created.ID, false))

	_, err = suite.postgres.GetBoardgameByID(ctx, created.ID)
	suite.Assert().Error(err)
}

func (suite *BoardgameSuite) TestDelete_NotFound() {
	ctx := context.Background()
	err := suite.postgres.DeleteBoardgame(ctx, 99999, false)
	suite.Assert().Error(err)
}

func (suite *BoardgameSuite) TestGetAll_FilterLike() {
	ctx := context.Background()
	_, err := suite.postgres.CreateBoardgame(ctx, testCreateInput("Catan"))
	suite.Require().NoError(err)
	_, err = suite.postgres.CreateBoardgame(ctx, testCreateInput("Vagrantsong"))
	suite.Require().NoError(err)

	all, err := suite.postgres.GetAllBoardgames(ctx, listopt.Params{Filter: listopt.Filter{Column: "name", Op: "like", Value: "cat"}}, false)
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 1)
	suite.Assert().Equal("Catan", all[0].Name)
}

func (suite *BoardgameSuite) TestGetAll_FilterNumeric() {
	ctx := context.Background()
	input1 := testCreateInput("Catan")
	input1.MaxPlayers = 4
	_, err := suite.postgres.CreateBoardgame(ctx, input1)
	suite.Require().NoError(err)

	input2 := testCreateInput("Vagrantsong")
	input2.MaxPlayers = 2
	_, err = suite.postgres.CreateBoardgame(ctx, input2)
	suite.Require().NoError(err)

	all, err := suite.postgres.GetAllBoardgames(ctx, listopt.Params{Filter: listopt.Filter{Column: "max_players", Op: "lt", Value: "3"}}, false)
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 1)
	suite.Assert().Equal("Vagrantsong", all[0].Name)
}

func TestBoardgameSuite(t *testing.T) {
	suite.Run(t, new(BoardgameSuite))
}
