package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/middleware"
)

type BoardGameSuite struct {
	suite.Suite
	base *Base
}

func (suite *BoardGameSuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *BoardGameSuite) TestPostBoardgameSuccess() {
	suite.base.dbMock.EXPECT().
		GetBoardgameBySlug(gomock.Any(), "test").
		Return(boardgame.Boardgame{}, middleware.NewError(http.StatusNotFound, "record not found"))
	suite.base.dbMock.EXPECT().
		CreateBoardgame(gomock.Any(), gomock.Any()).
		Return(boardgame.Boardgame{Slug: "test", Name: "test", MinPlayers: 1, MaxPlayers: 4}, nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame").
		JSON(`{"name":"test","min_players":1,"max_players":4}`).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *BoardGameSuite) TestPostExpansion() {
	parent := boardgame.Boardgame{ID: 1, Slug: "parent", Name: "parent", MinPlayers: 1, MaxPlayers: 4}
	suite.base.dbMock.EXPECT().
		GetBoardgameBySlug(gomock.Any(), "expansion").
		Return(boardgame.Boardgame{}, middleware.NewError(http.StatusNotFound, "record not found"))
	suite.base.dbMock.EXPECT().
		GetBoardgameBySlug(gomock.Any(), "parent").
		Return(parent, nil)

	suite.base.dbMock.EXPECT().
		CreateBoardgame(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, input boardgame.CreateBoardgameDTO) (boardgame.Boardgame, error) {
			suite.Require().NotNil(input.ParentID)
			suite.Equal(uint(1), *input.ParentID)
			return boardgame.Boardgame{
				Slug:        "expansion",
				Name:        "expansion",
				MinPlayers:  1,
				MaxPlayers:  4,
				BoardgameID: input.ParentID,
			}, nil
		})

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/boardgame/parent/expansion").
		JSON(`{"name":"expansion","min_players":1,"max_players":4}`).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *BoardGameSuite) TestGetBoardgame() {
	expected := boardgame.Boardgame{}
	suite.base.dbMock.EXPECT().
		GetBoardgameBySlug(gomock.Any(), "test").
		Return(expected, nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/boardgame/test").
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *BoardGameSuite) TestDeleteBoardgameSuccess() {
	bg := boardgame.Boardgame{ID: 1, Slug: "test"}
	suite.base.dbMock.EXPECT().
		GetBoardgameBySlug(gomock.Any(), "test").
		Return(bg, nil)
	suite.base.dbMock.EXPECT().
		DeleteBoardgame(gomock.Any(), uint(1), false).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/boardgame/test").
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()
}

func TestBoardGameSuite(t *testing.T) {
	suite.Run(t, new(BoardGameSuite))
}
