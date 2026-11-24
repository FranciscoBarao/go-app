package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BoardgameSuite struct {
	suite.Suite
}

func (suite *BoardgameSuite) TestCreate_Success() {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	db := &Postgres{pool: &fakeDB{
		tx:  &fakeTx{row: valueRow{values: []any{uint(3)}}},
		row: valueRow{values: boardgameValues(uint(3), "catan", "Catan", now)},
	}}

	created, err := db.CreateBoardgame(context.Background(), boardgame.CreateBoardgameDTO{
		Slug: "catan", Name: "Catan", MinPlayers: 1, MaxPlayers: 4,
	})

	suite.Require().NoError(err)
	suite.Assert().Equal(uint(3), created.ID)
	suite.Assert().Equal("Catan", created.Name)
	suite.Assert().Empty(created.Categories)
}

func (suite *BoardgameSuite) TestCreate_AlreadyExists() {
	db := &Postgres{pool: &fakeDB{tx: &fakeTx{row: errRow{err: pgError(uniqueViolationCode)}}}}

	_, err := db.CreateBoardgame(context.Background(), boardgame.CreateBoardgameDTO{Slug: "catan"})

	suite.Assert().ErrorIs(err, store.ErrAlreadyExists)
}

func (suite *BoardgameSuite) TestCreate_MissingCategory() {
	db := &Postgres{pool: &fakeDB{tx: &fakeTx{
		row:      valueRow{values: []any{uint(3)}},
		execErrs: []error{pgError(foreignKeyViolationCode)},
	}}}

	_, err := db.CreateBoardgame(context.Background(), boardgame.CreateBoardgameDTO{
		Slug: "catan", Categories: []uint{9},
	})

	suite.Assert().ErrorIs(err, store.ErrForeignKey)
}

func (suite *BoardgameSuite) TestCreate_BeginError() {
	db := &Postgres{pool: &fakeDB{beginErr: assert.AnError}}

	_, err := db.CreateBoardgame(context.Background(), boardgame.CreateBoardgameDTO{Slug: "catan"})

	suite.Assert().EqualError(err, "database operation failed")
	suite.Assert().NotErrorIs(err, assert.AnError)
}

func (suite *BoardgameSuite) TestGetBySlug_Success() {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	db := &Postgres{pool: &fakeDB{row: valueRow{values: boardgameValues(uint(3), "catan", "Catan", now)}}}

	record, err := db.GetBoardgameBySlug(context.Background(), "catan")

	suite.Require().NoError(err)
	suite.Assert().Equal("Catan", record.Name)
	suite.Assert().Empty(record.Mechanisms)
	suite.Assert().Empty(record.Expansions)
}

func (suite *BoardgameSuite) TestGetBySlug_NotFound() {
	db := &Postgres{pool: &fakeDB{row: errRow{err: pgx.ErrNoRows}}}

	_, err := db.GetBoardgameBySlug(context.Background(), "missing")

	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *BoardgameSuite) TestGetAll_Success() {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	db := &Postgres{pool: &fakeDB{
		row: valueRow{values: []any{1}},
		queued: []pgx.Rows{&fakeRows{scans: []func(dest ...any) error{
			func(dest ...any) error {
				return assign(dest, boardgameValues(uint(3), "catan", "Catan", now)...)
			},
		}}},
	}}

	records, total, err := db.GetAllBoardgames(context.Background(), listopt.Query{}, false)

	suite.Require().NoError(err)
	suite.Assert().Equal(1, total)
	suite.Require().Len(records, 1)
	suite.Assert().Equal("catan", records[0].Slug)
}

func (suite *BoardgameSuite) TestGetAll_CountError() {
	db := &Postgres{pool: &fakeDB{row: errRow{err: assert.AnError}}}

	_, _, err := db.GetAllBoardgames(context.Background(), listopt.Query{}, true)

	suite.Assert().EqualError(err, "database operation failed")
}

func (suite *BoardgameSuite) TestUpdate_Success() {
	db := &Postgres{pool: &fakeDB{}}

	err := db.UpdateBoardgame(context.Background(), &boardgame.Boardgame{ID: 3, Name: "Catan"})

	suite.Assert().NoError(err)
}

func (suite *BoardgameSuite) TestUpdate_NotFound() {
	db := &Postgres{pool: &fakeDB{execSet: true, execTag: updated("UPDATE 0")}}

	err := db.UpdateBoardgame(context.Background(), &boardgame.Boardgame{ID: 99})

	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func (suite *BoardgameSuite) TestDelete_Success() {
	db := &Postgres{pool: &fakeDB{tx: &fakeTx{}}}

	err := db.DeleteBoardgame(context.Background(), 3, false)

	suite.Assert().NoError(err)
}

func (suite *BoardgameSuite) TestDelete_NotFound() {
	db := &Postgres{pool: &fakeDB{tx: &fakeTx{execTags: []pgconn.CommandTag{
		updated("UPDATE 1"),
		updated("UPDATE 0"),
	}}}}

	err := db.DeleteBoardgame(context.Background(), 99, false)

	suite.Assert().ErrorIs(err, store.ErrNotFound)
}

func boardgameValues(id uint, slug, name string, now time.Time) []any {
	return []any{
		id, slug, now, now, (*time.Time)(nil),
		name, "", 1995, 1, 4, 0, 0,
		(*int)(nil), (*int)(nil), (*uint)(nil),
	}
}

func TestBoardgameSuite(t *testing.T) {
	suite.Run(t, new(BoardgameSuite))
}
