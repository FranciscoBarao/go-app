package tests

import (
	"fmt"
	"testing"

	"github.com/FranciscoBarao/catalog/config"
	"github.com/FranciscoBarao/catalog/database"
	"github.com/FranciscoBarao/catalog/model"

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
	suite.postgres, err = database.Connect(config)
	suite.Require().NoError(err)
}

func (suite *PostgresSuite) TearDownSuite() {
	suite.db.Shutdown()
}

func (suite *PostgresSuite) InsertEntry(model interface{}) {
	err := suite.postgres.Create(model)
	suite.Require().NoError(err)
}

func (suite *PostgresSuite) TestGetBoardgames() {
	insertBg := &model.Boardgame{Name: "name-01", Publisher: "publisher", PlayerNumber: 1}
	suite.InsertEntry(insertBg)

	insertBg = &model.Boardgame{Name: "name-02", Publisher: "publisher", PlayerNumber: 2}
	suite.InsertEntry(insertBg)

	var readBg []model.Boardgame
	err := suite.postgres.Read(&readBg, "", "", "")
	suite.Assert().NoError(err)
	suite.Assert().Len(readBg, 2)
}

func (suite *PostgresSuite) TestDeleteBoardgame() {
	name := "deleteName"
	insertBg := &model.Boardgame{Name: name, Publisher: "publisher", PlayerNumber: 1}
	suite.InsertEntry(insertBg)

	var readBg *model.Boardgame
	err := suite.postgres.Read(&readBg, "", "name = ?", name)
	suite.Assert().NoError(err)

	suite.postgres.Delete(readBg)
}

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}
