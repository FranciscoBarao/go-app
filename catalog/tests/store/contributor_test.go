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
	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/store"
	"github.com/FranciscoBarao/catalog/internal/store/postgres"
)

type ContributorSuite struct {
	suite.Suite
	db    *postgres.Postgres
	admin *pgxpool.Pool
}

func (suite *ContributorSuite) SetupSuite() {
	port, stop, err := startTestPostgres("catalog-contributor")
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

func (suite *ContributorSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Shutdown()
	}
	if suite.admin != nil {
		suite.admin.Close()
	}
}

func (suite *ContributorSuite) TearDownTest() {
	_, err := suite.admin.Exec(context.Background(), `TRUNCATE TABLE contributors RESTART IDENTITY CASCADE`)
	suite.Require().NoError(err)
}

func (suite *ContributorSuite) TestRoundTripUpdateAndErrors() {
	ctx := context.Background()
	bio := "designs games"
	created := &contributor.Contributor{Slug: "reiner-knizia", Name: "Reiner Knizia", Bio: &bio}
	suite.Require().NoError(suite.db.CreateContributor(ctx, created))

	got, err := suite.db.GetContributorBySlug(ctx, "reiner-knizia")
	suite.Assert().NoError(err)
	suite.Assert().Equal("Reiner Knizia", got.Name)
	suite.Require().NotNil(got.Bio)
	suite.Assert().Equal(bio, *got.Bio)

	id, err := suite.db.GetContributorIDBySlug(ctx, "reiner-knizia")
	suite.Assert().NoError(err)
	suite.Assert().Equal(created.ID, id)

	updatedBio := "still designing"
	got.Name = "Dr. Knizia"
	got.Bio = &updatedBio
	suite.Require().NoError(suite.db.UpdateContributor(ctx, &got))

	reloaded, err := suite.db.GetContributorBySlug(ctx, "reiner-knizia")
	suite.Assert().NoError(err)
	suite.Assert().Equal("Dr. Knizia", reloaded.Name)
	suite.Require().NotNil(reloaded.Bio)
	suite.Assert().Equal(updatedBio, *reloaded.Bio)
	suite.Assert().True(reloaded.UpdatedAt.After(created.UpdatedAt) || reloaded.UpdatedAt.Equal(created.UpdatedAt))

	missing := contributor.Contributor{Slug: "missing", Name: "Nobody"}
	err = suite.db.UpdateContributor(ctx, &missing)
	suite.Assert().ErrorIs(err, store.ErrNotFound)

	err = suite.db.CreateContributor(ctx, &contributor.Contributor{Slug: "reiner-knizia", Name: "Again"})
	suite.Assert().ErrorIs(err, store.ErrAlreadyExists)

	_, err = suite.db.GetContributorBySlug(ctx, "missing")
	suite.Assert().ErrorIs(err, store.ErrNotFound)

	suite.Require().NoError(suite.db.DeleteContributor(ctx, "reiner-knizia", false))
	_, err = suite.db.GetContributorBySlug(ctx, "reiner-knizia")
	suite.Assert().ErrorIs(err, store.ErrNotFound)

	err = suite.db.DeleteContributor(ctx, "missing", false)
	suite.Assert().ErrorIs(err, store.ErrNotFound)

	suite.Require().NoError(suite.db.CreateContributor(ctx, &contributor.Contributor{Slug: "kosmos", Name: "Kosmos"}))
	listed, total, err := suite.db.GetAllContributors(ctx, listopt.NewQuery(listopt.WithFilter(listopt.Filter{
		Column:   "name",
		Operator: listopt.Like,
		Value:    "kos",
	})))
	suite.Assert().NoError(err)
	suite.Assert().Equal(1, total)
	suite.Require().Len(listed, 1)
	suite.Assert().Equal("kosmos", listed[0].Slug)
}

func TestContributorSuite(t *testing.T) {
	suite.Run(t, new(ContributorSuite))
}
