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
	"github.com/FranciscoBarao/catalog/internal/middleware"
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
	req := boardgame.CreateBoardgameRequest{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4}
	expectedBg := boardgame.NewBoardgame(&req)

	suite.mockSvc.EXPECT().Create(gomock.Any(), expectedBg, uint(0)).Return(nil)

	bgBytes, err := json.Marshal(req)
	suite.Require().NoError(err)
	body := bytes.NewReader(bgBytes)

	httpReq := httptest.NewRequest(http.MethodPost, "/", body)
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, httpReq)

	suite.Assert().Equal(http.StatusOK, rec.Code)
	var result boardgame.Boardgame
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Assert().Equal(req.Name, result.Name)
	suite.Assert().Equal(req.Publisher, result.Publisher)
	suite.Assert().Equal(req.PlayerNumber, result.PlayerNumber)

}

func (suite *BoardgameControllerSuite) TestCreate_InternalError() {
	req := boardgame.CreateBoardgameRequest{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4}
	expectedBg := boardgame.NewBoardgame(&req)

	suite.mockSvc.EXPECT().
		Create(gomock.Any(), expectedBg, uint(0)).
		Return(assert.AnError)

	bgBytes, err := json.Marshal(req)
	suite.Require().NoError(err)
	body := bytes.NewReader(bgBytes)

	httpReq := httptest.NewRequest(http.MethodPost, "/", body)
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, httpReq)

	suite.Assert().Equal(http.StatusInternalServerError, rec.Code)
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
	suite.mockSvc.EXPECT().
		GetByID(gomock.Any(), uint(1)).
		Return(expected, nil)

	req := reqWithParam(httptest.NewRequest(http.MethodGet, "/", nil), "id", "1")
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result boardgame.Boardgame
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal(expected.Name, result.Name)
}

func (suite *BoardgameControllerSuite) TestGet_NotFound() {
	suite.mockSvc.EXPECT().
		GetByID(gomock.Any(), uint(999)).
		Return(boardgame.Boardgame{}, middleware.NewError(http.StatusNotFound, "not found"))

	req := reqWithParam(httptest.NewRequest(http.MethodGet, "/", nil), "id", "999")
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *BoardgameControllerSuite) TestGetAll() {
	expected := []boardgame.Boardgame{
		{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4},
		{Name: "Vagrantsong", Publisher: "Karma", PlayerNumber: 2},
	}
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), "").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result []boardgame.Boardgame
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Len(result, 2)
}

func (suite *BoardgameControllerSuite) TestGetAll_Empty() {
	expected := []boardgame.Boardgame{}
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), "").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result []boardgame.Boardgame
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Empty(result)
}

func (suite *BoardgameControllerSuite) TestGetAll_InternalError() {
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), "").Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *BoardgameControllerSuite) TestUpdate() {
	name := "Catan Updated"
	dto := boardgame.UpdateBoardgameRequest{Name: &name}

	suite.mockSvc.EXPECT().
		Update(gomock.Any(), &dto, uint(1)).
		Return(nil)

	bgBytes, err := json.Marshal(dto)
	suite.Require().NoError(err)
	body := bytes.NewReader(bgBytes)

	req := reqWithParam(
		httptest.NewRequest(http.MethodPatch, "/", body),
		"id", "1",
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Update(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *BoardgameControllerSuite) TestUpdate_InvalidInput() {
	playerNumber := 20 // Invalid input
	dto := boardgame.UpdateBoardgameRequest{PlayerNumber: &playerNumber}

	bgBytes, err := json.Marshal(dto)
	suite.Require().NoError(err)
	body := bytes.NewReader(bgBytes)

	req := reqWithParam(
		httptest.NewRequest(http.MethodPatch, "/", body),
		"id", "1",
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Update(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *BoardgameControllerSuite) TestUpdate_InvalidID() {
	name := "Catan"
	dto := boardgame.UpdateBoardgameRequest{Name: &name}

	bgBytes, err := json.Marshal(dto)
	suite.Require().NoError(err)
	body := bytes.NewReader(bgBytes)

	req := reqWithParam(
		httptest.NewRequest(http.MethodPatch, "/", body),
		"id", "invalid", // Invalid ID
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Update(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *BoardgameControllerSuite) TestUpdate_InternalError() {
	name := "Catan"
	dto := boardgame.UpdateBoardgameRequest{Name: &name}

	suite.mockSvc.EXPECT().
		Update(gomock.Any(), &dto, uint(1)).
		Return(assert.AnError)

	bgBytes, err := json.Marshal(dto)
	suite.Require().NoError(err)
	body := bytes.NewReader(bgBytes)

	req := reqWithParam(
		httptest.NewRequest(http.MethodPatch, "/", body),
		"id", "1",
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Update(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *BoardgameControllerSuite) TestUpdate_NotFound() {
	name := "Catan"
	dto := boardgame.UpdateBoardgameRequest{Name: &name}

	suite.mockSvc.EXPECT().
		Update(gomock.Any(), &dto, uint(1)).
		Return(middleware.NewError(http.StatusNotFound, "not found"))

	bgBytes, err := json.Marshal(dto)
	suite.Require().NoError(err)
	body := bytes.NewReader(bgBytes)

	req := reqWithParam(
		httptest.NewRequest(http.MethodPatch, "/", body),
		"id", "1",
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Update(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *BoardgameControllerSuite) TestDelete() {
	suite.mockSvc.EXPECT().
		DeleteByID(gomock.Any(), uint(1)).
		Return(nil)

	req := reqWithParam(httptest.NewRequest(http.MethodDelete, "/", nil), "id", "1")
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusNoContent, rec.Code)
}

func (suite *BoardgameControllerSuite) TestDelete_InvalidID() {
	req := reqWithParam(
		httptest.NewRequest(http.MethodDelete, "/", nil),
		"id", "invalid", // Invalid ID
	)
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *BoardgameControllerSuite) TestDelete_InternalError() {
	suite.mockSvc.EXPECT().
		DeleteByID(gomock.Any(), uint(1)).
		Return(assert.AnError)

	req := reqWithParam(
		httptest.NewRequest(http.MethodDelete, "/", nil),
		"id", "1",
	)
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *BoardgameControllerSuite) TestDelete_NotFound() {
	suite.mockSvc.EXPECT().
		DeleteByID(gomock.Any(), uint(1)).
		Return(middleware.NewError(http.StatusNotFound, "not found"))

	req := reqWithParam(
		httptest.NewRequest(http.MethodDelete, "/", nil),
		"id", "1",
	)
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *BoardgameControllerSuite) TestRate() {
	expectedRating := &boardgame.Rating{Username: "testuser", Value: 8}

	suite.mockSvc.EXPECT().
		Rate(gomock.Any(), gomock.Any(), uint(1), "testuser").
		Return(nil)

	bgBytes, err := json.Marshal(expectedRating)
	suite.Require().NoError(err)
	body := bytes.NewReader(bgBytes)

	req := reqWithParam(
		httptest.NewRequest(http.MethodPost, "/", body),
		"id", "1",
	)
	req.Header.Set("Content-Type", "application/json")
	claims := map[string]string{"username": "testuser"}
	req = req.WithContext(context.WithValue(req.Context(), oauth.ClaimsContext, claims))
	rec := httptest.NewRecorder()

	suite.controller.Rate(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *BoardgameControllerSuite) TestRate_InvalidInput() {
	expectedRating := &boardgame.Rating{
		Username: "testuser",
		Value:    11, // Invalid value
	}

	bgBytes, err := json.Marshal(expectedRating)
	suite.Require().NoError(err)
	body := bytes.NewReader(bgBytes)

	req := reqWithParam(
		httptest.NewRequest(http.MethodPost, "/", body),
		"id", "1",
	)
	req.Header.Set("Content-Type", "application/json")
	claims := map[string]string{"username": "testuser"}
	req = req.WithContext(context.WithValue(req.Context(), oauth.ClaimsContext, claims))
	rec := httptest.NewRecorder()

	suite.controller.Rate(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *BoardgameControllerSuite) TestRate_InvalidID() {
	expectedRating := &boardgame.Rating{Username: "testuser", Value: 7}

	bgBytes, err := json.Marshal(expectedRating)
	suite.Require().NoError(err)
	body := bytes.NewReader(bgBytes)

	req := reqWithParam(
		httptest.NewRequest(http.MethodPost, "/", body),
		"id", "invalid", // Invalid ID
	)
	req.Header.Set("Content-Type", "application/json")
	claims := map[string]string{"username": "testuser"}
	req = req.WithContext(context.WithValue(req.Context(), oauth.ClaimsContext, claims))
	rec := httptest.NewRecorder()

	suite.controller.Rate(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *BoardgameControllerSuite) TestRate_InternalError() {
	expectedRating := &boardgame.Rating{Username: "testuser", Value: 8}

	suite.mockSvc.EXPECT().
		Rate(gomock.Any(), gomock.Any(), uint(1), "testuser").
		Return(assert.AnError)

	bgBytes, err := json.Marshal(expectedRating)
	suite.Require().NoError(err)
	body := bytes.NewReader(bgBytes)

	req := reqWithParam(
		httptest.NewRequest(http.MethodPost, "/", body),
		"id", "1",
	)
	req.Header.Set("Content-Type", "application/json")
	claims := map[string]string{"username": "testuser"}
	req = req.WithContext(context.WithValue(req.Context(), oauth.ClaimsContext, claims))
	rec := httptest.NewRecorder()

	suite.controller.Rate(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *BoardgameControllerSuite) TestRate_NotFound() {
	expectedRating := &boardgame.Rating{Username: "testuser", Value: 8}

	suite.mockSvc.EXPECT().
		Rate(gomock.Any(), gomock.Any(), uint(1), "testuser").
		Return(middleware.NewError(http.StatusNotFound, "not found"))

	bgBytes, err := json.Marshal(expectedRating)
	suite.Require().NoError(err)
	body := bytes.NewReader(bgBytes)

	req := reqWithParam(
		httptest.NewRequest(http.MethodPost, "/", body),
		"id", "1",
	)
	req.Header.Set("Content-Type", "application/json")
	claims := map[string]string{"username": "testuser"}
	req = req.WithContext(context.WithValue(req.Context(), oauth.ClaimsContext, claims))
	rec := httptest.NewRecorder()

	suite.controller.Rate(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func TestBoardgameControllerSuite(t *testing.T) {
	suite.Run(t, new(BoardgameControllerSuite))
}
