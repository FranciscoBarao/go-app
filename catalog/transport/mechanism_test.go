package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/mechanism"
	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/stretchr/testify/assert"
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
	expectedMechanism := &mechanism.Mechanism{Name: "DeckBuilding"}
	suite.mockSvc.EXPECT().
		Create(gomock.Any(), expectedMechanism).
		Return(nil)

	mechanismBytes, err := json.Marshal(expectedMechanism)
	suite.Require().NoError(err)
	body := bytes.NewReader(mechanismBytes)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result mechanism.Mechanism
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal(expectedMechanism.Name, result.Name)
}

func (suite *MechanismControllerSuite) TestCreate_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *MechanismControllerSuite) TestCreate_InvalidStruct() {
	const maxChars = 30
	longName := strings.Repeat("a", maxChars+1)
	expectedmechanism := &mechanism.Mechanism{
		Name: longName, // invalid name length
	}

	mechanismBytes, err := json.Marshal(expectedmechanism)
	suite.Require().NoError(err)
	body := bytes.NewReader(mechanismBytes)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *MechanismControllerSuite) TestCreate_InternalError() {
	expectedmechanism := &mechanism.Mechanism{Name: "DeckBuilding"}

	suite.mockSvc.EXPECT().
		Create(gomock.Any(), expectedmechanism).
		Return(assert.AnError)

	mechanismBytes, err := json.Marshal(expectedmechanism)
	suite.Require().NoError(err)
	body := bytes.NewReader(mechanismBytes)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MechanismControllerSuite) TestGetAll() {
	expected := []mechanism.Mechanism{{Name: "DeckBuilding"}, {Name: "WorkerPlacement"}}
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), "").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result []mechanism.Mechanism
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Len(result, 2)
}

func (suite *MechanismControllerSuite) TestGetAll_Empty() {
	expected := []mechanism.Mechanism{}
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), "").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result []mechanism.Mechanism
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Empty(result)
}

func (suite *MechanismControllerSuite) TestGetAll_InternalError() {
	suite.mockSvc.EXPECT().
		GetAll(gomock.Any(), "").Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MechanismControllerSuite) TestGet() {
	name := "DeckBuilding"
	suite.mockSvc.EXPECT().
		Get(gomock.Any(), name).
		Return(mechanism.Mechanism{Name: name}, nil)

	req := reqWithParam(
		httptest.NewRequest(http.MethodGet, "/", nil),
		"name", name,
	)
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result mechanism.Mechanism
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal("DeckBuilding", result.Name)
}

func (suite *MechanismControllerSuite) TestGet_NotFound() {
	name := "not found"
	suite.mockSvc.EXPECT().
		Get(gomock.Any(), name).
		Return(mechanism.Mechanism{}, middleware.NewError(http.StatusNotFound, "not found"))

	req := reqWithParam(
		httptest.NewRequest(http.MethodGet, "/", nil),
		"name", name,
	)
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *MechanismControllerSuite) TestGet_InternalError() {
	name := "error"
	suite.mockSvc.EXPECT().
		Get(gomock.Any(), name).
		Return(mechanism.Mechanism{}, assert.AnError)

	req := reqWithParam(
		httptest.NewRequest(http.MethodGet, "/", nil),
		"name", name,
	)
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MechanismControllerSuite) TestDelete() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "DeckBuilding").Return(nil)

	req := reqWithParam(httptest.NewRequest(http.MethodDelete, "/", nil), "name", "DeckBuilding")
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusNoContent, rec.Code)
}

func (suite *MechanismControllerSuite) TestDelete_NotFound() {
	name := "not found"
	suite.mockSvc.EXPECT().
		Delete(gomock.Any(), name).
		Return(middleware.NewError(http.StatusNotFound, "not found"))

	req := reqWithParam(
		httptest.NewRequest(http.MethodDelete, "/", nil),
		"name", name,
	)
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *MechanismControllerSuite) TestDelete_InternalError() {
	name := "Strategy"
	suite.mockSvc.EXPECT().
		Delete(gomock.Any(), name).
		Return(assert.AnError)

	req := reqWithParam(
		httptest.NewRequest(http.MethodDelete, "/", nil),
		"name", name,
	)
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestMechanismControllerSuite(t *testing.T) {
	suite.Run(t, new(MechanismControllerSuite))
}
