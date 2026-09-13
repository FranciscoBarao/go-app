package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/listopt"
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

func (suite *BoardGameSuite) TestGetBoardgamesWithQueryParams() {
	expected := []boardgame.Boardgame{{Slug: "catan", Name: "Catan", MinPlayers: 3, MaxPlayers: 4}}
	var gotQuery listopt.Query
	suite.base.dbMock.EXPECT().
		GetAllBoardgames(gomock.Any(), gomock.Any(), true).
		DoAndReturn(func(_ context.Context, query listopt.Query, _ bool) ([]boardgame.Boardgame, int, error) {
			gotQuery = query
			return expected, 12, nil
		})

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/boardgame").
		QueryParams(map[string]string{
			"page":            "2",
			"pageSize":        "5",
			"sort":            "name.asc",
			"filter":          "min_players.ge.3",
			"include_deleted": "true",
		}).
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(assertEnvelope(suite.T(), 1, 2, 5, 12, 3)).
		End()

	suite.Assert().Equal("name", gotQuery.Sort.Column)
	suite.Assert().Equal("asc", gotQuery.Sort.Order)
	suite.Require().Len(gotQuery.Filters, 1)
	suite.Assert().Equal("min_players", gotQuery.Filters[0].Column)
	suite.Assert().Equal(listopt.Ge, gotQuery.Filters[0].Operator)
	suite.Assert().Equal(5, gotQuery.Pagination.Limit())
	suite.Assert().Equal(5, gotQuery.Pagination.Offset())
}

func (suite *BoardGameSuite) TestGetBoardgamesMalformedQueryParam() {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/boardgame").
		QueryParams(map[string]string{"filter": "min_players.bogus.3"}).
		Expect(suite.T()).
		Status(http.StatusUnprocessableEntity).
		End()
}

func (suite *BoardGameSuite) TestQueryBoardgames() {
	expected := []boardgame.Boardgame{{Slug: "catan", Name: "Catan", MinPlayers: 3, MaxPlayers: 4}}
	suite.base.dbMock.EXPECT().
		GetAllBoardgames(gomock.Any(), gomock.Any(), false).
		Return(expected, len(expected), nil)

	assertEnvelope := func(res *http.Response, _ *http.Request) error {
		var body struct {
			Data       []boardgame.Boardgame `json:"data"`
			Page       int                   `json:"page"`
			PageSize   int                   `json:"pageSize"`
			TotalItems int                   `json:"totalItems"`
			TotalPages int                   `json:"totalPages"`
		}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return err
		}
		suite.Assert().Len(body.Data, 1)
		suite.Assert().Equal(1, body.Page)
		suite.Assert().Equal(10, body.PageSize)
		suite.Assert().Equal(1, body.TotalItems)
		suite.Assert().Equal(1, body.TotalPages)
		return nil
	}

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Method("QUERY").
		URL("/api/boardgame").
		Body(`{"pagination":{"page":1,"pageSize":10},"filters":[{"field":"min_players","op":"ge","value":"3"}]}`).
		ContentType("application/json").
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(assertEnvelope).
		End()
}

func TestBoardGameSuite(t *testing.T) {
	suite.Run(t, new(BoardGameSuite))
}
