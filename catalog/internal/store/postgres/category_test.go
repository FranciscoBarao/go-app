package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CategorySuite struct {
	suite.Suite
}

func (suite *CategorySuite) TestCreate_Success() {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	record := &category.Category{Slug: "strategy", Name: "Strategy"}
	db := &Postgres{pool: &fakeDB{row: valueRow{values: []any{uint(7), now, now}}}}

	err := db.CreateCategory(context.Background(), record)

	suite.Require().NoError(err)
	suite.Assert().Equal(uint(7), record.ID)
	suite.Assert().Equal(now, record.CreatedAt)
	suite.Assert().Equal(now, record.UpdatedAt)
}

func (suite *CategorySuite) TestCreate_AlreadyExists() {
	db := &Postgres{pool: &fakeDB{row: errRow{err: pgError(uniqueViolationCode)}}}

	err := db.CreateCategory(context.Background(), &category.Category{Slug: "strategy"})

	suite.Assert().ErrorIs(err, store.ErrAlreadyExists)
}

func (suite *CategorySuite) TestCreate_ForeignKey() {
	db := &Postgres{pool: &fakeDB{row: errRow{err: pgError(foreignKeyViolationCode)}}}

	err := db.CreateCategory(context.Background(), &category.Category{Slug: "strategy"})

	suite.Assert().ErrorIs(err, store.ErrForeignKey)
}

func (suite *CategorySuite) TestCreate_UnexpectedError() {
	db := &Postgres{pool: &fakeDB{row: errRow{err: assert.AnError}}}

	err := db.CreateCategory(context.Background(), &category.Category{Slug: "strategy"})

	suite.Assert().EqualError(err, "database operation failed")
	suite.Assert().NotErrorIs(err, assert.AnError)
}

func (suite *CategorySuite) TestGetBySlug_Success() {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	db := &Postgres{pool: &fakeDB{queued: []pgx.Rows{categoryRows(uint(7), "strategy", "Strategy", now)}}}

	record, err := db.GetCategoryBySlug(context.Background(), "strategy")

	suite.Require().NoError(err)
	suite.Assert().Equal(uint(7), record.ID)
	suite.Assert().Equal("Strategy", record.Name)
}

func (suite *CategorySuite) TestGetBySlug_NotFound() {
	db := &Postgres{pool: &fakeDB{}}

	_, err := db.GetCategoryBySlug(context.Background(), "missing")

	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *CategorySuite) TestGetBySlug_QueryError() {
	db := &Postgres{pool: &fakeDB{queryErr: assert.AnError}}

	_, err := db.GetCategoryBySlug(context.Background(), "strategy")

	suite.Assert().EqualError(err, "database operation failed")
}

func (suite *CategorySuite) TestGetIDBySlug_Success() {
	db := &Postgres{pool: &fakeDB{row: valueRow{values: []any{uint(4)}}}}

	id, err := db.GetCategoryIDBySlug(context.Background(), "strategy")

	suite.Require().NoError(err)
	suite.Assert().Equal(uint(4), id)
}

func (suite *CategorySuite) TestGetIDBySlug_NotFound() {
	db := &Postgres{pool: &fakeDB{row: errRow{err: pgx.ErrNoRows}}}

	_, err := db.GetCategoryIDBySlug(context.Background(), "missing")

	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *CategorySuite) TestDelete_Success() {
	db := &Postgres{pool: &fakeDB{}}

	err := db.DeleteCategory(context.Background(), "strategy", false)

	suite.Assert().NoError(err)
}

func (suite *CategorySuite) TestDelete_NotFound() {
	db := &Postgres{pool: &fakeDB{execSet: true, execTag: updated("UPDATE 0")}}

	err := db.DeleteCategory(context.Background(), "missing", false)

	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *CategorySuite) TestDelete_ForeignKey() {
	db := &Postgres{pool: &fakeDB{execErr: pgError(foreignKeyViolationCode)}}

	err := db.DeleteCategory(context.Background(), "strategy", true)

	suite.Assert().ErrorIs(err, store.ErrForeignKey)
}

func (suite *CategorySuite) TestGetAll_Success() {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	db := &Postgres{pool: &fakeDB{
		row:    valueRow{values: []any{1}},
		queued: []pgx.Rows{categoryRows(uint(7), "strategy", "Strategy", now)},
	}}

	records, total, err := db.GetAllCategories(context.Background(), listopt.Query{})

	suite.Require().NoError(err)
	suite.Assert().Equal(1, total)
	suite.Require().Len(records, 1)
	suite.Assert().Equal("strategy", records[0].Slug)
}

func (suite *CategorySuite) TestGetAll_CountError() {
	db := &Postgres{pool: &fakeDB{row: errRow{err: assert.AnError}}}

	_, _, err := db.GetAllCategories(context.Background(), listopt.Query{})

	suite.Assert().EqualError(err, "database operation failed")
}

func categoryRows(id uint, slug, name string, now time.Time) *fakeRows {
	return &fakeRows{scans: []func(dest ...any) error{
		func(dest ...any) error {
			return assign(dest, id, slug, name, (*int)(nil), now, now, (*time.Time)(nil))
		},
	}}
}

func TestCategorySuite(t *testing.T) {
	suite.Run(t, new(CategorySuite))
}
