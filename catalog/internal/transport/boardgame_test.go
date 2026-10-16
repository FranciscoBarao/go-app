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
	"github.com/go-chi/oauth"
	"github.com/stretchr/testify/assert"
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

func (suite *BoardgameControllerSuite) TestCreate_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (suite *BoardgameControllerSuite) TestCreate_InvalidPlayers() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Catan","min_players":0,"max_players":4}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (suite *BoardgameControllerSuite) TestCreate_AlreadyExists() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any(), "").Return(boardgame.Boardgame{}, boardgame.ErrAlreadyExists)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Catan","min_players":2,"max_players":4}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Assert().Equal(http.StatusConflict, rec.Code)
}

func (suite *BoardgameControllerSuite) TestCreate_NotFound() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any(), "missing").Return(boardgame.Boardgame{}, boardgame.ErrNotFound)

	req := withRouteParam(httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Scenario","min_players":2,"max_players":4}`)), "slug", "missing")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Assert().Equal(http.StatusNotFound, rec.Code)
}

func (suite *BoardgameControllerSuite) TestCreate_Validation() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any(), "").Return(boardgame.Boardgame{}, boardgame.ErrValidation)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Catan","min_players":5,"max_players":2}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (suite *BoardgameControllerSuite) TestCreate_MissingReference() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any(), "").Return(boardgame.Boardgame{}, boardgame.ErrMissingReference)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Catan","min_players":2,"max_players":4}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Assert().Equal(http.StatusConflict, rec.Code)
}

func (suite *BoardgameControllerSuite) TestCreate_ParentIsExpansion() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any(), "seafarers").Return(boardgame.Boardgame{}, boardgame.ErrParentIsExpansion)

	req := withRouteParam(httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Scenario","min_players":2,"max_players":4}`)), "slug", "seafarers")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Assert().Equal(http.StatusConflict, rec.Code)
}

func (suite *BoardgameControllerSuite) TestCreate_StoreError() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any(), "").Return(boardgame.Boardgame{}, assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Catan","min_players":2,"max_players":4}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Assert().Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *BoardgameControllerSuite) TestGet() {
	suite.mockSvc.EXPECT().GetBySlug(gomock.Any(), "catan").Return(boardgame.Boardgame{Slug: "catan"}, nil)

	httpReq := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "slug", "catan")
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, httpReq)
	suite.Assert().Equal(http.StatusOK, rec.Code)
}

func (suite *BoardgameControllerSuite) TestGet_NotFound() {
	suite.mockSvc.EXPECT().GetBySlug(gomock.Any(), "missing").Return(boardgame.Boardgame{}, boardgame.ErrNotFound)

	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "slug", "missing")
	rec := httptest.NewRecorder()
	suite.controller.Get(rec, req)
	suite.Assert().Equal(http.StatusNotFound, rec.Code)
}

func (suite *BoardgameControllerSuite) TestGet_StoreError() {
	suite.mockSvc.EXPECT().GetBySlug(gomock.Any(), "catan").Return(boardgame.Boardgame{}, assert.AnError)

	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "slug", "catan")
	rec := httptest.NewRecorder()
	suite.controller.Get(rec, req)
	suite.Assert().Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *BoardgameControllerSuite) TestGetByID() {
	suite.mockSvc.EXPECT().GetByID(gomock.Any(), uint(9)).Return(boardgame.Boardgame{ID: 9, Slug: "catan"}, nil)

	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "id", "9")
	rec := httptest.NewRecorder()
	suite.controller.GetByID(rec, req)
	suite.Assert().Equal(http.StatusOK, rec.Code)
}

func (suite *BoardgameControllerSuite) TestGetByID_Invalid() {
	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "id", "nope")
	rec := httptest.NewRecorder()
	suite.controller.GetByID(rec, req)
	suite.Assert().Equal(http.StatusBadRequest, rec.Code)
	suite.Assert().Contains(rec.Body.String(), "invalid boardgame id nope")
}

func (suite *BoardgameControllerSuite) TestGetByID_NotFound() {
	suite.mockSvc.EXPECT().GetByID(gomock.Any(), uint(9)).Return(boardgame.Boardgame{}, boardgame.ErrNotFound)

	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "id", "9")
	rec := httptest.NewRecorder()
	suite.controller.GetByID(rec, req)
	suite.Assert().Equal(http.StatusNotFound, rec.Code)
}

func (suite *BoardgameControllerSuite) TestGetByID_StoreError() {
	suite.mockSvc.EXPECT().GetByID(gomock.Any(), uint(9)).Return(boardgame.Boardgame{}, assert.AnError)

	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "id", "9")
	rec := httptest.NewRecorder()
	suite.controller.GetByID(rec, req)
	suite.Assert().Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *BoardgameControllerSuite) TestGetAll() {
	expected := []boardgame.Boardgame{{Slug: "catan"}}
	var gotQuery listopt.Query
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), true, gomock.Any()).DoAndReturn(
		func(_ context.Context, _ bool, query listopt.Query) ([]boardgame.Boardgame, int, error) {
			gotQuery = query
			return expected, 42, nil
		},
	)

	httpReq := httptest.NewRequest(http.MethodGet, "/?page=2&pageSize=20&sort=name.desc&filter=min_players.ge.3&include_deleted=true", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, httpReq)
	suite.Require().Equal(http.StatusOK, rec.Code)

	suite.Assert().Equal(2, gotQuery.Pagination.Page)
	suite.Assert().Equal(20, gotQuery.Pagination.PageSize)
	suite.Assert().Equal("name", gotQuery.Sort.Column)
	suite.Assert().Equal("desc", gotQuery.Sort.Order)
	suite.Require().Len(gotQuery.Filters, 1)
	suite.Assert().Equal("min_players", gotQuery.Filters[0].Column)

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

func (suite *BoardgameControllerSuite) TestGetAll_StoreError() {
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), false, gomock.Any()).Return(nil, 0, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	suite.controller.GetAll(rec, req)
	suite.Assert().Equal(http.StatusInternalServerError, rec.Code)
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

func (suite *BoardgameControllerSuite) TestUpdate() {
	suite.mockSvc.EXPECT().Update(gomock.Any(), gomock.Any(), "catan").Return(nil)

	req := withRouteParam(httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"description":"A trading game"}`)), "slug", "catan")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Update(rec, req)
	suite.Assert().Equal(http.StatusOK, rec.Code)
}

func (suite *BoardgameControllerSuite) TestUpdate_NotFound() {
	suite.mockSvc.EXPECT().Update(gomock.Any(), gomock.Any(), "missing").Return(boardgame.ErrNotFound)

	req := withRouteParam(httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"description":"A trading game"}`)), "slug", "missing")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Update(rec, req)
	suite.Assert().Equal(http.StatusNotFound, rec.Code)
}

func (suite *BoardgameControllerSuite) TestUpdate_InvalidPlayers() {
	req := withRouteParam(httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"min_players":99}`)), "slug", "catan")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Update(rec, req)
	suite.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (suite *BoardgameControllerSuite) TestUpdate_Conflict() {
	suite.mockSvc.EXPECT().Update(gomock.Any(), gomock.Any(), "catan").Return(boardgame.ErrAlreadyExists)

	req := withRouteParam(httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"description":"A trading game"}`)), "slug", "catan")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Update(rec, req)
	suite.Assert().Equal(http.StatusConflict, rec.Code)
}

func (suite *BoardgameControllerSuite) TestUpdate_StoreError() {
	suite.mockSvc.EXPECT().Update(gomock.Any(), gomock.Any(), "catan").Return(assert.AnError)

	req := withRouteParam(httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"description":"A trading game"}`)), "slug", "catan")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Update(rec, req)
	suite.Assert().Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *BoardgameControllerSuite) TestDelete() {
	suite.mockSvc.EXPECT().DeleteBySlug(gomock.Any(), "catan", true).Return(nil)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/?hard=true", nil), "slug", "catan")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Assert().Equal(http.StatusNoContent, rec.Code)
}

func (suite *BoardgameControllerSuite) TestDelete_NotFound() {
	suite.mockSvc.EXPECT().DeleteBySlug(gomock.Any(), "missing", false).Return(boardgame.ErrNotFound)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/", nil), "slug", "missing")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Assert().Equal(http.StatusNotFound, rec.Code)
}

func (suite *BoardgameControllerSuite) TestDelete_StoreError() {
	suite.mockSvc.EXPECT().DeleteBySlug(gomock.Any(), "catan", false).Return(assert.AnError)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/", nil), "slug", "catan")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Assert().Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *BoardgameControllerSuite) TestRate() {
	suite.mockSvc.EXPECT().Rate(gomock.Any(), gomock.Any(), "catan", "ada").DoAndReturn(
		func(_ context.Context, rating *boardgame.Rating, _ string, username string) error {
			suite.Equal(8, rating.Value)
			rating.Username = username
			return nil
		},
	)

	req := withRouteParam(httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"username":"ada","value":8}`)), "slug", "catan")
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), oauth.ClaimsContext, map[string]string{"username": "ada"}))
	rec := httptest.NewRecorder()
	suite.controller.Rate(rec, req)
	suite.Assert().Equal(http.StatusOK, rec.Code)
	suite.Assert().Contains(rec.Body.String(), `"username":"ada"`)
}

func (suite *BoardgameControllerSuite) TestRate_Unauthorized() {
	req := withRouteParam(httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"username":"ada","value":8}`)), "slug", "catan")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Rate(rec, req)
	suite.Assert().Equal(http.StatusUnauthorized, rec.Code)
}

func (suite *BoardgameControllerSuite) TestRate_MissingUsernameClaim() {
	req := withRouteParam(httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"username":"ada","value":8}`)), "slug", "catan")
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), oauth.ClaimsContext, map[string]string{}))
	rec := httptest.NewRecorder()
	suite.controller.Rate(rec, req)
	suite.Assert().Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *BoardgameControllerSuite) TestRate_InvalidValue() {
	req := withRouteParam(httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"username":"ada","value":11}`)), "slug", "catan")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Rate(rec, req)
	suite.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (suite *BoardgameControllerSuite) TestRate_NotFound() {
	suite.mockSvc.EXPECT().Rate(gomock.Any(), gomock.Any(), "missing", "ada").Return(boardgame.ErrNotFound)

	req := withRouteParam(httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"username":"ada","value":8}`)), "slug", "missing")
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), oauth.ClaimsContext, map[string]string{"username": "ada"}))
	rec := httptest.NewRecorder()
	suite.controller.Rate(rec, req)
	suite.Assert().Equal(http.StatusNotFound, rec.Code)
}

func TestBoardgameControllerSuite(t *testing.T) {
	suite.Run(t, new(BoardgameControllerSuite))
}
