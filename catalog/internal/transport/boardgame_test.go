package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type BoardgameControllerSuite struct {
	suite.Suite
	mockSvc    *MockBoardgameService
	controller *BoardgameController
}

func (suite *BoardgameControllerSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockSvc = NewMockBoardgameService(ctrl)
	suite.controller = NewBoardgameController(suite.mockSvc)
}

func (suite *BoardgameControllerSuite) TestCreate() {
	req := boardgame.CreateBoardgameRequest{Name: "Catan", MinPlayers: 2, MaxPlayers: 4}
	expected := boardgame.Boardgame{Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}

	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any(), "").DoAndReturn(
		func(_ context.Context, r *boardgame.CreateBoardgameRequest, _ string) (boardgame.Boardgame, error) {
			suite.Equal(req.Name, r.Name)
			suite.Equal(req.MinPlayers, r.MinPlayers)
			return expected, nil
		},
	)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, httpReq)
	suite.Assert().Equal(http.StatusOK, rec.Code)
}

func (suite *BoardgameControllerSuite) TestGet() {
	suite.mockSvc.EXPECT().GetBySlug(gomock.Any(), "catan").Return(boardgame.Boardgame{Slug: "catan"}, nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", "catan")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, httpReq)
	suite.Assert().Equal(http.StatusOK, rec.Code)
}

func (suite *BoardgameControllerSuite) TestGetAll() {
	expected := []boardgame.Boardgame{{Slug: "catan"}}
	var gotParams listopt.Params
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), true, gomock.Any()).DoAndReturn(
		func(_ context.Context, _ bool, opts ...listopt.Option) ([]boardgame.Boardgame, int, error) {
			gotParams = listopt.Apply(opts...)
			return expected, 42, nil
		},
	)

	httpReq := httptest.NewRequest(http.MethodGet, "/?page=2&pageSize=20&sort=name.desc&filter=minplayers.ge.3&include_deleted=true", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, httpReq)
	suite.Require().Equal(http.StatusOK, rec.Code)

	suite.Assert().Equal(2, gotParams.Pagination.Page)
	suite.Assert().Equal(20, gotParams.Pagination.PageSize)
	suite.Assert().Equal("name", gotParams.Sort.Column)
	suite.Assert().Equal("desc", gotParams.Sort.Order)
	suite.Require().Len(gotParams.Filters, 1)
	suite.Assert().Equal("min_players", gotParams.Filters[0].Column)

	var resp PaginatedResponse[boardgame.Boardgame]
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &resp))
	suite.Assert().Equal(2, resp.Page)
	suite.Assert().Equal(20, resp.PageSize)
	suite.Assert().Equal(42, resp.TotalItems)
	suite.Assert().Equal(3, resp.TotalPages)
	suite.Assert().Len(resp.Data, 1)
}

func (suite *BoardgameControllerSuite) TestGetAll_NoParamsUsesDefaults() {
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), false, gomock.Any()).Return([]boardgame.Boardgame{}, 0, nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, httpReq)
	suite.Require().Equal(http.StatusOK, rec.Code)

	var resp PaginatedResponse[boardgame.Boardgame]
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &resp))
	suite.Assert().Equal(listopt.DefaultPage, resp.Page)
	suite.Assert().Equal(listopt.DefaultPageSize, resp.PageSize)
	suite.Assert().Equal(0, resp.TotalItems)
	suite.Assert().NotNil(resp.Data)
}

func (suite *BoardgameControllerSuite) TestGetAll_MalformedParams() {
	for _, rawQuery := range []string{"page=abc", "pageSize=abc", "filter=name.cat", "sort=name", "filter=name.bogus.x", "sort=categories.asc"} {
		suite.Run(rawQuery, func() {
			httpReq := httptest.NewRequest(http.MethodGet, "/?"+rawQuery, nil)
			rec := httptest.NewRecorder()

			suite.controller.GetAll(rec, httpReq)
			suite.Assert().Equal(http.StatusUnprocessableEntity, rec.Code)
		})
	}
}

func (suite *BoardgameControllerSuite) TestQuery() {
	expected := []boardgame.Boardgame{{Slug: "catan"}}
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), false, gomock.Any()).Return(expected, 1, nil)

	httpReq := httptest.NewRequest("QUERY", "/", strings.NewReader(`{"pagination":{"page":1,"pageSize":10}}`))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, httpReq)
	suite.Assert().Equal(http.StatusOK, rec.Code)

	var resp PaginatedResponse[boardgame.Boardgame]
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &resp))
	suite.Assert().Equal(1, resp.TotalItems)
	suite.Assert().Equal(1, resp.TotalPages)
	suite.Assert().Len(resp.Data, 1)
}

func (suite *BoardgameControllerSuite) TestQuery_IncludeDeletedFlowsToService() {
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), true, gomock.Any()).Return([]boardgame.Boardgame{}, 0, nil)

	httpReq := httptest.NewRequest("QUERY", "/", strings.NewReader(`{"include_deleted":true,"pagination":{"page":1,"pageSize":10}}`))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, httpReq)
	suite.Assert().Equal(http.StatusOK, rec.Code)
}

func (suite *BoardgameControllerSuite) TestQuery_ClampsPaginationInEnvelope() {
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), false, gomock.Any()).Return([]boardgame.Boardgame{}, 0, nil)

	httpReq := httptest.NewRequest("QUERY", "/", strings.NewReader(`{"pagination":{"page":0,"pageSize":500}}`))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, httpReq)
	suite.Assert().Equal(http.StatusOK, rec.Code)

	var resp PaginatedResponse[boardgame.Boardgame]
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &resp))
	suite.Assert().Equal(1, resp.Page)
	suite.Assert().Equal(100, resp.PageSize)
}

func (suite *BoardgameControllerSuite) TestQuery_EmptyObjectUsesDefaults() {
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), false, gomock.Any()).Return([]boardgame.Boardgame{}, 0, nil)

	httpReq := httptest.NewRequest("QUERY", "/", strings.NewReader(`{}`))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, httpReq)
	suite.Assert().Equal(http.StatusOK, rec.Code)

	var resp PaginatedResponse[boardgame.Boardgame]
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &resp))
	suite.Assert().Equal(1, resp.Page)
	suite.Assert().Equal(10, resp.PageSize)
	suite.Assert().Equal(0, resp.TotalItems)
	suite.Assert().NotNil(resp.Data)
	suite.Assert().Len(resp.Data, 0)
}

func (suite *BoardgameControllerSuite) TestQuery_EmptyBodyReturnsBadRequest() {
	httpReq := httptest.NewRequest("QUERY", "/", strings.NewReader(""))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, httpReq)
	suite.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (suite *BoardgameControllerSuite) TestQuery_MissingContentTypeReturnsBadRequest() {
	httpReq := httptest.NewRequest("QUERY", "/", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, httpReq)
	suite.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (suite *BoardgameControllerSuite) TestQuery_InvalidFilterOp() {
	httpReq := httptest.NewRequest("QUERY", "/", strings.NewReader(`{"filters":[{"field":"name","op":"bogus","value":"x"}]}`))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, httpReq)
	suite.Assert().Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (suite *BoardgameControllerSuite) TestQuery_InvalidSortField() {
	httpReq := httptest.NewRequest("QUERY", "/", strings.NewReader(`{"sort":{"field":"categories","order":"asc"}}`))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, httpReq)
	suite.Assert().Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (suite *BoardgameControllerSuite) TestCreate_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func TestBoardgameControllerSuite(t *testing.T) {
	suite.Run(t, new(BoardgameControllerSuite))
}
