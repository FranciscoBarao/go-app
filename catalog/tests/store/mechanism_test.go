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
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/store"
	"github.com/FranciscoBarao/catalog/internal/store/postgres"
)

type MechanismSuite struct {
	suite.Suite
	db    *postgres.Postgres
	admin *pgxpool.Pool
}

func (suite *MechanismSuite) SetupSuite() {
	port, stop, err := startTestPostgres("catalog-mechanism")
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

func (suite *MechanismSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Shutdown()
	}
	if suite.admin != nil {
		suite.admin.Close()
	}
}

func (suite *MechanismSuite) TearDownTest() {
	_, err := suite.admin.Exec(context.Background(), `TRUNCATE TABLE mechanisms RESTART IDENTITY CASCADE`)
	suite.Require().NoError(err)
}

func (suite *MechanismSuite) TestRoundTripAndErrors() {
	ctx := context.Background()
	created := &mechanism.Mechanism{Slug: "worker-placement", Name: "Worker Placement"}
	suite.Require().NoError(suite.db.CreateMechanism(ctx, created))

	got, err := suite.db.GetMechanismBySlug(ctx, "worker-placement")
	suite.Assert().NoError(err)
	suite.Assert().Equal(created.ID, got.ID)

	id, err := suite.db.GetMechanismIDBySlug(ctx, "worker-placement")
	suite.Assert().NoError(err)
	suite.Assert().Equal(created.ID, id)

	err = suite.db.CreateMechanism(ctx, &mechanism.Mechanism{Slug: "worker-placement", Name: "Again"})
	suite.Assert().ErrorIs(err, store.ErrAlreadyExists)

	_, err = suite.db.GetMechanismBySlug(ctx, "missing")
	suite.Assert().ErrorIs(err, store.ErrNotFound)

	suite.Require().NoError(suite.db.DeleteMechanism(ctx, "worker-placement", false))
	_, err = suite.db.GetMechanismBySlug(ctx, "worker-placement")
	suite.Assert().ErrorIs(err, store.ErrNotFound)

	err = suite.db.DeleteMechanism(ctx, "missing", true)
	suite.Assert().ErrorIs(err, store.ErrNotFound)

	err = suite.db.CreateMechanism(ctx, &mechanism.Mechanism{Slug: "auction", Name: "Auction"})
	suite.Require().NoError(err)
	listed, total, err := suite.db.GetAllMechanisms(ctx, listopt.NewQuery(listopt.WithFilter(listopt.Filter{
		Column:   "name",
		Operator: listopt.Like,
		Value:    "auc",
	})))
	suite.Assert().NoError(err)
	suite.Assert().Equal(1, total)
	suite.Require().Len(listed, 1)
	suite.Assert().Equal("auction", listed[0].Slug)
}

func TestMechanismSuite(t *testing.T) {
	suite.Run(t, new(MechanismSuite))
}
