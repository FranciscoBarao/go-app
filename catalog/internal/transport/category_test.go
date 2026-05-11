package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/stretchr/testify/assert"
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
	expectedCategory := &category.Category{Name: "Strategy"}

	suite.mockSvc.EXPECT().
		Create(gomock.Any(), expectedCategory).
		Return(nil)

	categoryBytes, err := json.Marshal(expectedCategory)
	suite.Require().NoError(err)
	body := bytes.NewReader(categoryBytes)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result category.Category
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal(expectedCategory.Name, result.Name)
}

func (suite *CategoryControllerSuite) TestCreate_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *CategoryControllerSuite) TestCreate_InvalidStruct() {
	const maxChars = 30
	longName := strings.Repeat("a", maxChars+1)
	expectedCategory := &category.Category{
		Name: longName, // invalid name length
	}

	categoryBytes, err := json.Marshal(expectedCategory)
	suite.Require().NoError(err)
	body := bytes.NewReader(categoryBytes)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *CategoryControllerSuite) TestCreate_InternalError() {
	expectedCategory := &category.Category{Name: "DeckBuilding"}

	suite.mockSvc.EXPECT().
		Create(gomock.Any(), expectedCategory).
		Return(assert.AnError)

	categoryBytes, err := json.Marshal(expectedCategory)
	suite.Require().NoError(err)
	body := bytes.NewReader(categoryBytes)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CategoryControllerSuite) TestGetAll() {
	expected := []category.Category{{Name: "Strategy"}, {Name: "Family"}}
	suite.mockSvc.EXPECT().GetAll(gomock.Any()).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result []category.Category
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Len(result, 2)
}

func (suite *CategoryControllerSuite) TestGetAll_Empty() {
	expected := []category.Category{}
	suite.mockSvc.EXPECT().GetAll(gomock.Any()).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result []category.Category
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Empty(result)
}

func (suite *CategoryControllerSuite) TestGetAll_InternalError() {
	suite.mockSvc.EXPECT().
		GetAll(gomock.Any()).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CategoryControllerSuite) TestGet() {
	suite.mockSvc.EXPECT().Get(gomock.Any(), "Strategy").Return(category.Category{Name: "Strategy"}, nil)

	req := reqWithParam(httptest.NewRequest(http.MethodGet, "/", nil), "name", "Strategy")
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result category.Category
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal("Strategy", result.Name)
}

func (suite *CategoryControllerSuite) TestGet_NotFound() {
	name := "not found"
	suite.mockSvc.EXPECT().
		Get(gomock.Any(), name).
		Return(category.Category{}, middleware.NewError(http.StatusNotFound, "not found"))

	req := reqWithParam(
		httptest.NewRequest(http.MethodGet, "/", nil),
		"name", name,
	)
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *CategoryControllerSuite) TestGet_InternalError() {
	name := "error"
	suite.mockSvc.EXPECT().
		Get(gomock.Any(), name).
		Return(category.Category{}, assert.AnError)

	req := reqWithParam(
		httptest.NewRequest(http.MethodGet, "/", nil),
		"name", name,
	)
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CategoryControllerSuite) TestDelete() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "Strategy").Return(nil)

	req := reqWithParam(httptest.NewRequest(http.MethodDelete, "/", nil), "name", "Strategy")
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusNoContent, rec.Code)
}

func (suite *CategoryControllerSuite) TestDelete_NotFound() {
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

func (suite *CategoryControllerSuite) TestDelete_InternalError() {
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

func TestCategoryControllerSuite(t *testing.T) {
	suite.Run(t, new(CategoryControllerSuite))
}
