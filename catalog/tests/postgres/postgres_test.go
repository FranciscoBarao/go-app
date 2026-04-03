package tests

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/FranciscoBarao/catalog/boardgame"
	"github.com/FranciscoBarao/catalog/category"
	"github.com/FranciscoBarao/catalog/config"
	"github.com/FranciscoBarao/catalog/database"
	"github.com/FranciscoBarao/catalog/mechanism"
	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/tag"

	"github.com/stretchr/testify/suite"
)

type PostgresSuite struct {
	suite.Suite
	db       *TestPostgres
	postgres *database.Postgres
}

func (suite *PostgresSuite) SetupSuite() {
	var err error
	suite.db, err = NewTestPostgres("test")
	suite.Require().NoError(err)

	// Fetch DB configs
	config := &config.PostgresConfig{
		Host:     "localhost",
		Username: "postgres",
		Password: "postgres",
		Port:     fmt.Sprintf("%d", suite.db.Port),
		Database: "postgres",
	}
	// Connect to Database
	suite.postgres, err = database.Connect(config, &boardgame.Boardgame{}, &tag.Tag{}, &category.Category{}, &mechanism.Mechanism{}, &boardgame.Rating{})
	suite.Require().NoError(err)
}

func (suite *PostgresSuite) TearDownSuite() {
	suite.db.Shutdown()
}

func (suite *PostgresSuite) InsertEntry(model interface{}) {
	err := suite.postgres.Create(model)
	suite.Require().NoError(err)
}

func (suite *PostgresSuite) TestCreate_IncorrectModel() {
	incorrectModel := &map[string]string{
		"a": "a",
	}
	err := suite.postgres.Create(incorrectModel)
	suite.Require().Error(err)
}

func (suite *PostgresSuite) TestCreate_Duplicated() {
	insertBg := &boardgame.Boardgame{Name: "duplicated", Publisher: "publisher", PlayerNumber: 1}
	err := suite.postgres.Create(insertBg)
	suite.Require().NoError(err)
	err = suite.postgres.Create(insertBg)
	suite.Require().Error(err)
}

func (suite *PostgresSuite) TestGet() {
	insertBg := &boardgame.Boardgame{Name: "name-01", Publisher: "publisher", PlayerNumber: 1}
	suite.InsertEntry(insertBg)

	insertBg = &boardgame.Boardgame{Name: "name-02", Publisher: "publisher", PlayerNumber: 2}
	suite.InsertEntry(insertBg)

	var readBg []boardgame.Boardgame
	err := suite.postgres.Read(&readBg, "", "name LIKE ?", "name-0%")
	suite.Assert().NoError(err)
	suite.Assert().Len(readBg, 2)
}

func (suite *PostgresSuite) TestUpdate() {
	name := "updateName"
	insertBg := &boardgame.Boardgame{Name: name, Publisher: "pub1", PlayerNumber: 1}
	suite.InsertEntry(insertBg)

	var bg boardgame.Boardgame
	err := suite.postgres.Read(&bg, "", "name = ?", name)
	suite.Assert().NoError(err)
	suite.Assert().Equal(insertBg.Publisher, bg.Publisher)

	newPublisher := "newpub2"
	bg.Publisher = newPublisher
	err = suite.postgres.Update(&bg)
	suite.Assert().NoError(err)

	err = suite.postgres.Read(&bg, "", "name = ?", name)
	suite.Assert().NoError(err)
	suite.Assert().Equal(newPublisher, bg.Publisher)
}

func (suite *PostgresSuite) TestDelete() {
	name := "deleteName"
	insertBg := &boardgame.Boardgame{Name: name, Publisher: "publisher", PlayerNumber: 1}
	suite.InsertEntry(insertBg)

	var readBg *boardgame.Boardgame
	err := suite.postgres.Read(&readBg, "", "name = ?", name)
	suite.Assert().NoError(err)

	err = suite.postgres.Delete(readBg)
	suite.Assert().NoError(err)
}

func (suite *PostgresSuite) TestDelete_NotFound() {
	var bg = &boardgame.Boardgame{}
	bg.ID = 100
	err := suite.postgres.Delete(bg)
	suite.Assert().Error(err)

	var mr = &middleware.MalformedRequest{}
	suite.Assert().True(errors.As(err, &mr))
	suite.Assert().Equal(mr.GetStatus(), http.StatusNotFound)
}

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}
