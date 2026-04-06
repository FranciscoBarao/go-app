package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/mechanism"
	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type MechanismControllerSuite struct {
	suite.Suite
	mockSvc    *MockMechanismService
	controller *MechanismController
}

func (suite *MechanismControllerSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockSvc = NewMockMechanismService(ctrl)
	suite.controller = NewMechanismController(suite.mockSvc)
}

func (suite *MechanismControllerSuite) TestCreate() {
	suite.mockSvc.EXPECT().Create(gomock.Any()).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"DeckBuilding"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result mechanism.Mechanism
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal("DeckBuilding", result.Name)
}

func (suite *MechanismControllerSuite) TestGetAll() {
	expected := []mechanism.Mechanism{{Name: "DeckBuilding"}, {Name: "WorkerPlacement"}}
	suite.mockSvc.EXPECT().GetAll("").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result []mechanism.Mechanism
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Len(result, 2)
}

func (suite *MechanismControllerSuite) TestGet() {
	suite.mockSvc.EXPECT().Get("DeckBuilding").Return(mechanism.Mechanism{Name: "DeckBuilding"}, nil)

	req := reqWithParam(httptest.NewRequest(http.MethodGet, "/", nil), "name", "DeckBuilding")
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result mechanism.Mechanism
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal("DeckBuilding", result.Name)
}

func (suite *MechanismControllerSuite) TestDelete() {
	suite.mockSvc.EXPECT().Delete("DeckBuilding").Return(nil)

	req := reqWithParam(httptest.NewRequest(http.MethodDelete, "/", nil), "name", "DeckBuilding")
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusNoContent, rec.Code)
}

func (suite *MechanismControllerSuite) TestCreateInvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *MechanismControllerSuite) TestServiceError() {
	suite.mockSvc.EXPECT().Get("Unknown").Return(mechanism.Mechanism{}, middleware.NewError(http.StatusNotFound, "not found"))

	req := reqWithParam(httptest.NewRequest(http.MethodGet, "/", nil), "name", "Unknown")
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func TestMechanismControllerSuite(t *testing.T) {
	suite.Run(t, new(MechanismControllerSuite))
}
