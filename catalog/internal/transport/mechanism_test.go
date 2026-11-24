package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/mechanism"
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
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mechanism.Mechanism{Slug: "trading", Name: "Trading"}, nil)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"name":"Trading"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *MechanismControllerSuite) TestCreate_EmptyName() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *MechanismControllerSuite) TestCreate_AlreadyExists() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mechanism.Mechanism{}, mechanism.ErrAlreadyExists)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Trading"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Equal(http.StatusConflict, rec.Code)
}

func (suite *MechanismControllerSuite) TestCreate_Validation() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mechanism.Mechanism{}, mechanism.ErrValidation)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"!!!"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *MechanismControllerSuite) TestCreate_StoreError() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mechanism.Mechanism{}, assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Trading"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MechanismControllerSuite) TestGet() {
	expected := mechanism.Mechanism{Slug: "trading", Name: "Trading"}
	suite.mockSvc.EXPECT().Get(gomock.Any(), "trading").Return(expected, nil)

	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "slug", "trading")
	rec := httptest.NewRecorder()
	suite.controller.Get(rec, req)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *MechanismControllerSuite) TestGet_NotFound() {
	suite.mockSvc.EXPECT().Get(gomock.Any(), "missing").Return(mechanism.Mechanism{}, mechanism.ErrNotFound)

	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "slug", "missing")
	rec := httptest.NewRecorder()
	suite.controller.Get(rec, req)
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *MechanismControllerSuite) TestGet_StoreError() {
	suite.mockSvc.EXPECT().Get(gomock.Any(), "trading").Return(mechanism.Mechanism{}, assert.AnError)

	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "slug", "trading")
	rec := httptest.NewRecorder()
	suite.controller.Get(rec, req)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MechanismControllerSuite) TestGetAll_MalformedPage() {
	req := httptest.NewRequest(http.MethodGet, "/?page=nope", nil)
	rec := httptest.NewRecorder()
	suite.controller.GetAll(rec, req)
	suite.Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (suite *MechanismControllerSuite) TestGetAll_StoreError() {
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(nil, 0, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	suite.controller.GetAll(rec, req)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *MechanismControllerSuite) TestQuery() {
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return([]mechanism.Mechanism{{Slug: "trading"}}, 1, nil)

	req := httptest.NewRequest("QUERY", "/", bytes.NewReader([]byte(`{"pagination":{"page":1,"pageSize":10}}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, req)
	suite.Equal(http.StatusOK, rec.Code)

	var resp PaginatedResponse[mechanism.Mechanism]
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &resp))
	suite.Equal(1, resp.TotalItems)
	suite.Equal(1, resp.Page)
	suite.Equal(10, resp.PageSize)
	suite.Len(resp.Data, 1)
}

func (suite *MechanismControllerSuite) TestQuery_InvalidJSON() {
	req := httptest.NewRequest("QUERY", "/", strings.NewReader(`{`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Query(rec, req)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *MechanismControllerSuite) TestDelete() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "trading", true).Return(nil)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/?hard=true", nil), "slug", "trading")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Equal(http.StatusNoContent, rec.Code)
}

func (suite *MechanismControllerSuite) TestDelete_NotFound() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "missing", false).Return(mechanism.ErrNotFound)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/", nil), "slug", "missing")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *MechanismControllerSuite) TestDelete_StillReferenced() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "trading", false).Return(mechanism.ErrStillReferenced)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/", nil), "slug", "trading")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Equal(http.StatusConflict, rec.Code)
}

func (suite *MechanismControllerSuite) TestDelete_StoreError() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "trading", false).Return(assert.AnError)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/", nil), "slug", "trading")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestMechanismControllerSuite(t *testing.T) {
	suite.Run(t, new(MechanismControllerSuite))
}
