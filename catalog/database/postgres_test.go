package database

import (
	"catalog/config"
	"catalog/model"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PostgresSuite struct {
	suite.Suite
	db       *TestPostgres
	postgres *Postgres
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
	suite.postgres, err = Connect(config)
	suite.Require().NoError(err)
}

func (suite *PostgresSuite) TearDownSuite() {
	suite.db.Shutdown()
}

func (suite *PostgresSuite) TestGetBoardgame() {
	insertBg := model.NewBoardgame("name", "publsiher", 1, []model.Tag{}, []model.Category{}, []model.Mechanism{})

	var omits = []string{"Tags.*", "Categories.*", "Mechanisms.*", "Ratings.*"}

	err := suite.postgres.Create(insertBg, omits...)
	suite.Require().NoError(err)

	var readBg []model.Boardgame
	err = suite.postgres.Read(&readBg, "", "", "")
	suite.Assert().NoError(err)
	suite.Assert().Len(readBg, 1)
}

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}
