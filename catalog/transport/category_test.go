package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/category"
	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CategoryControllerSuite struct {
	suite.Suite
	mockSvc    *MockCategoryService
	controller *CategoryController
}

func (suite *CategoryControllerSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockSvc = NewMockCategoryService(ctrl)
	suite.controller = NewCategoryController(suite.mockSvc)
}

func (suite *CategoryControllerSuite) TestCreate() {
	suite.mockSvc.EXPECT().Create(gomock.Any()).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Strategy"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result category.Category
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal("Strategy", result.Name)
}

func (suite *CategoryControllerSuite) TestGetAll() {
	expected := []category.Category{{Name: "Strategy"}, {Name: "Family"}}
	suite.mockSvc.EXPECT().GetAll("").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result []category.Category
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Len(result, 2)
}

func (suite *CategoryControllerSuite) TestGet() {
	suite.mockSvc.EXPECT().Get("Strategy").Return(category.Category{Name: "Strategy"}, nil)

	req := reqWithParam(httptest.NewRequest(http.MethodGet, "/", nil), "name", "Strategy")
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result category.Category
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal("Strategy", result.Name)
}

func (suite *CategoryControllerSuite) TestDelete() {
	suite.mockSvc.EXPECT().Delete("Strategy").Return(nil)

	req := reqWithParam(httptest.NewRequest(http.MethodDelete, "/", nil), "name", "Strategy")
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusNoContent, rec.Code)
}

func (suite *CategoryControllerSuite) TestCreateInvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *CategoryControllerSuite) TestServiceError() {
	suite.mockSvc.EXPECT().Get("Unknown").Return(category.Category{}, middleware.NewError(http.StatusNotFound, "not found"))

	req := reqWithParam(httptest.NewRequest(http.MethodGet, "/", nil), "name", "Unknown")
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func TestCategoryControllerSuite(t *testing.T) {
	suite.Run(t, new(CategoryControllerSuite))
}
