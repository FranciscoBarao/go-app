package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MechanismSuite struct {
	suite.Suite
}

func (suite *MechanismSuite) TestCreate_Success() {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	record := &mechanism.Mechanism{Slug: "drafting", Name: "Drafting"}
	db := &Postgres{pool: &fakeDB{row: valueRow{values: []any{uint(3), now, now}}}}

	err := db.CreateMechanism(context.Background(), record)

	suite.Require().NoError(err)
	suite.Assert().Equal(uint(3), record.ID)
	suite.Assert().Equal(now, record.UpdatedAt)
}

func (suite *MechanismSuite) TestCreate_AlreadyExists() {
	db := &Postgres{pool: &fakeDB{row: errRow{err: pgError(uniqueViolationCode)}}}

	err := db.CreateMechanism(context.Background(), &mechanism.Mechanism{Slug: "drafting"})

	suite.Assert().ErrorIs(err, store.ErrAlreadyExists)
}

func (suite *MechanismSuite) TestCreate_UnexpectedError() {
	db := &Postgres{pool: &fakeDB{row: errRow{err: assert.AnError}}}

	err := db.CreateMechanism(context.Background(), &mechanism.Mechanism{Slug: "drafting"})

	suite.Assert().EqualError(err, "database operation failed")
	suite.Assert().NotErrorIs(err, assert.AnError)
}

func (suite *MechanismSuite) TestGetBySlug_Success() {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	db := &Postgres{pool: &fakeDB{queued: []pgx.Rows{mechanismRows(uint(3), "drafting", "Drafting", now)}}}

	record, err := db.GetMechanismBySlug(context.Background(), "drafting")

	suite.Require().NoError(err)
	suite.Assert().Equal("Drafting", record.Name)
}

func (suite *MechanismSuite) TestGetBySlug_NotFound() {
	db := &Postgres{pool: &fakeDB{}}

	_, err := db.GetMechanismBySlug(context.Background(), "missing")

	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *MechanismSuite) TestGetIDBySlug_NotFound() {
	db := &Postgres{pool: &fakeDB{row: errRow{err: pgx.ErrNoRows}}}

	_, err := db.GetMechanismIDBySlug(context.Background(), "missing")

	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *MechanismSuite) TestDelete_NotFound() {
	db := &Postgres{pool: &fakeDB{execSet: true, execTag: updated("DELETE 0")}}

	err := db.DeleteMechanism(context.Background(), "missing", true)

	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *MechanismSuite) TestDelete_ForeignKey() {
	db := &Postgres{pool: &fakeDB{execErr: pgError(foreignKeyViolationCode)}}

	err := db.DeleteMechanism(context.Background(), "drafting", false)

	suite.Assert().ErrorIs(err, store.ErrForeignKey)
}

func (suite *MechanismSuite) TestGetAll_QueryError() {
	db := &Postgres{pool: &fakeDB{
		row:      valueRow{values: []any{0}},
		queryErr: assert.AnError,
	}}

	_, _, err := db.GetAllMechanisms(context.Background(), listopt.Query{})

	suite.Assert().EqualError(err, "database operation failed")
}

func mechanismRows(id uint, slug, name string, now time.Time) *fakeRows {
	return &fakeRows{scans: []func(dest ...any) error{
		func(dest ...any) error {
			return assign(dest, id, slug, name, (*int)(nil), now, now, (*time.Time)(nil))
		},
	}}
}

func TestMechanismSuite(t *testing.T) {
	suite.Run(t, new(MechanismSuite))
}
