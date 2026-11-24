package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ContributorSuite struct {
	suite.Suite
}

func (suite *ContributorSuite) TestCreate_Success() {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	record := &contributor.Contributor{Slug: "knizia", Name: "Reiner Knizia"}
	db := &Postgres{pool: &fakeDB{row: valueRow{values: []any{uint(9), now, now}}}}

	err := db.CreateContributor(context.Background(), record)

	suite.Require().NoError(err)
	suite.Assert().Equal(uint(9), record.ID)
	suite.Assert().Equal(now, record.CreatedAt)
}

func (suite *ContributorSuite) TestCreate_AlreadyExists() {
	db := &Postgres{pool: &fakeDB{row: errRow{err: pgError(uniqueViolationCode)}}}

	err := db.CreateContributor(context.Background(), &contributor.Contributor{Slug: "knizia"})

	suite.Assert().ErrorIs(err, store.ErrAlreadyExists)
}

func (suite *ContributorSuite) TestUpdate_Success() {
	now := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	record := &contributor.Contributor{Slug: "knizia", Name: "Reiner Knizia"}
	db := &Postgres{pool: &fakeDB{row: valueRow{values: []any{now}}}}

	err := db.UpdateContributor(context.Background(), record)

	suite.Require().NoError(err)
	suite.Assert().Equal(now, record.UpdatedAt)
}

func (suite *ContributorSuite) TestUpdate_NotFound() {
	db := &Postgres{pool: &fakeDB{row: errRow{err: pgx.ErrNoRows}}}

	err := db.UpdateContributor(context.Background(), &contributor.Contributor{Slug: "missing"})

	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *ContributorSuite) TestGetBySlug_Success() {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	db := &Postgres{pool: &fakeDB{queued: []pgx.Rows{contributorRows(uint(9), "knizia", "Reiner Knizia", now)}}}

	record, err := db.GetContributorBySlug(context.Background(), "knizia")

	suite.Require().NoError(err)
	suite.Assert().Equal("Reiner Knizia", record.Name)
	suite.Assert().Nil(record.Bio)
}

func (suite *ContributorSuite) TestGetBySlug_NotFound() {
	db := &Postgres{pool: &fakeDB{}}

	_, err := db.GetContributorBySlug(context.Background(), "missing")

	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *ContributorSuite) TestGetIDBySlug_Success() {
	db := &Postgres{pool: &fakeDB{row: valueRow{values: []any{uint(9)}}}}

	id, err := db.GetContributorIDBySlug(context.Background(), "knizia")

	suite.Require().NoError(err)
	suite.Assert().Equal(uint(9), id)
}

func (suite *ContributorSuite) TestDelete_NotFound() {
	db := &Postgres{pool: &fakeDB{execSet: true, execTag: updated("UPDATE 0")}}

	err := db.DeleteContributor(context.Background(), "missing", false)

	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *ContributorSuite) TestDelete_UnexpectedError() {
	db := &Postgres{pool: &fakeDB{execErr: assert.AnError}}

	err := db.DeleteContributor(context.Background(), "knizia", false)

	suite.Assert().EqualError(err, "database operation failed")
	suite.Assert().NotErrorIs(err, assert.AnError)
}

func (suite *ContributorSuite) TestGetAll_Success() {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	db := &Postgres{pool: &fakeDB{
		row:    valueRow{values: []any{1}},
		queued: []pgx.Rows{contributorRows(uint(9), "knizia", "Reiner Knizia", now)},
	}}

	records, total, err := db.GetAllContributors(context.Background(), listopt.Query{})

	suite.Require().NoError(err)
	suite.Assert().Equal(1, total)
	suite.Require().Len(records, 1)
	suite.Assert().Equal("knizia", records[0].Slug)
}

func contributorRows(id uint, slug, name string, now time.Time) *fakeRows {
	return &fakeRows{scans: []func(dest ...any) error{
		func(dest ...any) error {
			return assign(dest, id, slug, name, (*string)(nil), (*int)(nil), now, now, (*time.Time)(nil))
		},
	}}
}

func TestContributorSuite(t *testing.T) {
	suite.Run(t, new(ContributorSuite))
}
