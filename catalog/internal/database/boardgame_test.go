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
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/listopt"

	"github.com/stretchr/testify/suite"
)

type BoardgameSuite struct {
	suite.Suite
	db       *TestPostgres
	postgres *Postgres
}

// migrationDir returns the absolute path to the migrations directory
// relative to this test file's location.
func migrationDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "migrations")
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
	const database = "boardgames"

	// Cascade will delete revisions since they are FKs
	_, err := suite.postgres.pool.Exec(context.Background(),
		fmt.Sprintf(`TRUNCATE TABLE %s CASCADE;`, database))
	suite.Require().NoError(err)
}

func (suite *BoardgameSuite) TestCreate_DuplicatedName() {
	ctx := context.Background()
	bg := &boardgame.Boardgame{Name: "duplicated", Publisher: "publisher", PlayerNumber: 1}
	err := suite.postgres.CreateBoardgame(ctx, bg)
	suite.Require().NoError(err)

	bg2 := &boardgame.Boardgame{Name: "duplicated", Publisher: "publisher", PlayerNumber: 1}
	err = suite.postgres.CreateBoardgame(ctx, bg2)
	suite.Assert().Error(err)
}

func (suite *BoardgameSuite) TestGet() {
	ctx := context.Background()
	expectedBg := &boardgame.Boardgame{Name: "name01", Publisher: "publisher", PlayerNumber: 1}
	err := suite.postgres.CreateBoardgame(ctx, expectedBg)
	suite.Require().NoError(err)

	bg, err := suite.postgres.GetBoardgameByID(ctx, expectedBg.ID)
	suite.Assert().NoError(err)
	suite.Assert().Equal(expectedBg.Name, bg.Name)
}

func (suite *BoardgameSuite) TestGetAll() {
	ctx := context.Background()
	bg1 := &boardgame.Boardgame{Name: "name01", Publisher: "publisher", PlayerNumber: 1}
	err := suite.postgres.CreateBoardgame(ctx, bg1)
	suite.Require().NoError(err)

	bg2 := &boardgame.Boardgame{Name: "name02", Publisher: "publisher", PlayerNumber: 2}
	err = suite.postgres.CreateBoardgame(ctx, bg2)
	suite.Require().NoError(err)

	all, err := suite.postgres.GetAllBoardgames(ctx, listopt.Params{Sort: listopt.Sort{Column: "id", Order: "asc"}})
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 2)
}

func (suite *BoardgameSuite) TestUpdate() {
	ctx := context.Background()
	name := "updateName"
	bg := &boardgame.Boardgame{Name: name, Publisher: "pub1", PlayerNumber: 1}
	err := suite.postgres.CreateBoardgame(ctx, bg)
	suite.Require().NoError(err)

	readBg, err := suite.postgres.GetBoardgameByID(ctx, bg.ID)
	suite.Assert().NoError(err)
	suite.Assert().Equal("pub1", readBg.Publisher)

	newPublisher := "newpub2"
	readBg.Publisher = newPublisher
	err = suite.postgres.UpdateBoardgame(ctx, &readBg)
	suite.Assert().NoError(err)

	updated, err := suite.postgres.GetBoardgameByID(ctx, bg.ID)
	suite.Assert().NoError(err)
	suite.Assert().Equal(newPublisher, updated.Publisher)
}

func (suite *BoardgameSuite) TestUpdate_NotFound() {
	ctx := context.Background()
	err := suite.postgres.UpdateBoardgame(ctx, &boardgame.Boardgame{Name: "not found"})
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().True(errors.As(err, &mr))
	suite.Assert().Equal(http.StatusNotFound, mr.GetStatus())
}

func (suite *BoardgameSuite) TestDelete() {
	ctx := context.Background()
	name := "deleteName"
	bg := &boardgame.Boardgame{Name: name, Publisher: "publisher", PlayerNumber: 1}
	err := suite.postgres.CreateBoardgame(ctx, bg)
	suite.Require().NoError(err)

	err = suite.postgres.DeleteBoardgame(ctx, bg.ID)
	suite.Assert().NoError(err)

	_, err = suite.postgres.GetBoardgameByID(ctx, bg.ID)
	suite.Assert().Error(err)
}

func (suite *BoardgameSuite) TestDelete_NotFound() {
	ctx := context.Background()
	err := suite.postgres.DeleteBoardgame(ctx, 99999)
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().True(errors.As(err, &mr))
	suite.Assert().Equal(http.StatusNotFound, mr.GetStatus())
}

func TestBoardgameSuite(t *testing.T) {
	suite.Run(t, new(BoardgameSuite))
}

func (suite *BoardgameSuite) TestGetAll_FilterLike() {
	ctx := context.Background()
	suite.Require().NoError(suite.postgres.CreateBoardgame(ctx, &boardgame.Boardgame{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4}))
	suite.Require().NoError(suite.postgres.CreateBoardgame(ctx, &boardgame.Boardgame{Name: "Vagrantsong", Publisher: "Wyrd", PlayerNumber: 2}))

	all, err := suite.postgres.GetAllBoardgames(ctx, listopt.Params{Filter: listopt.Filter{Column: "name", Op: "like", Value: "cat"}})
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 1)
	suite.Assert().Equal("Catan", all[0].Name)
}

func (suite *BoardgameSuite) TestGetAll_FilterNumeric() {
	ctx := context.Background()
	suite.Require().NoError(suite.postgres.CreateBoardgame(ctx, &boardgame.Boardgame{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4}))
	suite.Require().NoError(suite.postgres.CreateBoardgame(ctx, &boardgame.Boardgame{Name: "Vagrantsong", Publisher: "Wyrd", PlayerNumber: 2}))

	all, err := suite.postgres.GetAllBoardgames(ctx, listopt.Params{Filter: listopt.Filter{Column: "player_number", Op: "lt", Value: "3"}})
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 1)
	suite.Assert().Equal("Vagrantsong", all[0].Name)
}

func (suite *BoardgameSuite) TestGetAll_FilterEq() {
	ctx := context.Background()
	suite.Require().NoError(suite.postgres.CreateBoardgame(ctx, &boardgame.Boardgame{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4}))
	suite.Require().NoError(suite.postgres.CreateBoardgame(ctx, &boardgame.Boardgame{Name: "Vagrantsong", Publisher: "Wyrd", PlayerNumber: 2}))

	all, err := suite.postgres.GetAllBoardgames(ctx, listopt.Params{Filter: listopt.Filter{Column: "name", Op: listopt.OpEq, Value: "Catan"}})
	suite.Assert().NoError(err)
	suite.Assert().Len(all, 1)
	suite.Assert().Equal("Catan", all[0].Name)
}
