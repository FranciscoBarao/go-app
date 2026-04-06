package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/boardgame"
	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/go-chi/oauth"
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
	suite.mockSvc.EXPECT().Create(gomock.Any(), "").Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Catan","publisher":"Kosmos","playerNumber":4}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result boardgame.Boardgame
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal("Catan", result.Name)
}

func (suite *BoardgameControllerSuite) TestCreate_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *BoardgameControllerSuite) TestGet() {
	expected := boardgame.Boardgame{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4}
	suite.mockSvc.EXPECT().GetByID("1").Return(expected, nil)

	req := reqWithParam(httptest.NewRequest(http.MethodGet, "/", nil), "id", "1")
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result boardgame.Boardgame
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal("Catan", result.Name)
}

func (suite *BoardgameControllerSuite) TestGet_NotFound() {
	id := "999"
	suite.mockSvc.EXPECT().GetByID(id).Return(boardgame.Boardgame{}, middleware.NewError(http.StatusNotFound, "not found"))

	req := reqWithParam(httptest.NewRequest(http.MethodGet, "/", nil), "id", id)
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *BoardgameControllerSuite) TestGetAll() {
	expected := []boardgame.Boardgame{
		{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4},
	}
	suite.mockSvc.EXPECT().GetAll("", "", "").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result []boardgame.Boardgame
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Len(result, 1)
}

func (suite *BoardgameControllerSuite) TestUpdate() {
	suite.mockSvc.EXPECT().Update(gomock.Any(), "1").Return(nil)

	req := reqWithParam(
		httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"name":"CatanV2","publisher":"Kosmos","playerNumber":6}`)),
		"id", "1",
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Update(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result boardgame.Boardgame
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal("CatanV2", result.Name)
}

func (suite *BoardgameControllerSuite) TestDelete() {
	suite.mockSvc.EXPECT().DeleteByID("1").Return(nil)

	req := reqWithParam(httptest.NewRequest(http.MethodDelete, "/", nil), "id", "1")
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusNoContent, rec.Code)
}

func (suite *BoardgameControllerSuite) TestRate() {
	suite.mockSvc.EXPECT().Rate(gomock.Any(), "1", "testuser").Return(nil)

	req := reqWithParam(
		httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"value":8}`)),
		"id", "1",
	)
	req.Header.Set("Content-Type", "application/json")
	// Inject OAuth claims so GetUsernameFromToken works
	claims := map[string]string{"username": "testuser"}
	req = req.WithContext(context.WithValue(req.Context(), oauth.ClaimsContext, claims))
	rec := httptest.NewRecorder()

	suite.controller.Rate(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
}

func TestBoardgameControllerSuite(t *testing.T) {
	suite.Run(t, new(BoardgameControllerSuite))
}
