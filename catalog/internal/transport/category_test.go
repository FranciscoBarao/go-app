package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/listopt"
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
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(category.Category{Slug: "strategy", Name: "Strategy"}, nil)

	body := bytes.NewReader([]byte(`{"name":"Strategy"}`))
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *CategoryControllerSuite) TestCreate_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *CategoryControllerSuite) TestCreate_EmptyName() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)
	suite.Equal(http.StatusBadRequest, rec.Code)
	suite.Contains(rec.Body.String(), "model validation failed")
}

func (suite *CategoryControllerSuite) TestCreate_AlreadyExists() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(category.Category{}, category.ErrAlreadyExists)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Strategy"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)
	suite.Equal(http.StatusConflict, rec.Code)
}

func (suite *CategoryControllerSuite) TestCreate_Validation() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(category.Category{}, category.ErrValidation)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"!!!"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *CategoryControllerSuite) TestCreate_StoreError() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(category.Category{}, assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Strategy"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CategoryControllerSuite) TestGet() {
	expected := category.Category{Slug: "strategy", Name: "Strategy"}
	suite.mockSvc.EXPECT().Get(gomock.Any(), "strategy").Return(expected, nil)

	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "slug", "strategy")
	rec := httptest.NewRecorder()
	suite.controller.Get(rec, req)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *CategoryControllerSuite) TestGet_NotFound() {
	suite.mockSvc.EXPECT().Get(gomock.Any(), "missing").Return(category.Category{}, category.ErrNotFound)

	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "slug", "missing")
	rec := httptest.NewRecorder()
	suite.controller.Get(rec, req)
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *CategoryControllerSuite) TestGet_StoreError() {
	suite.mockSvc.EXPECT().Get(gomock.Any(), "strategy").Return(category.Category{}, assert.AnError)

	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "slug", "strategy")
	rec := httptest.NewRecorder()
	suite.controller.Get(rec, req)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CategoryControllerSuite) TestGetAll_MalformedPage() {
	req := httptest.NewRequest(http.MethodGet, "/?page=nope", nil)
	rec := httptest.NewRecorder()
	suite.controller.GetAll(rec, req)
	suite.Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (suite *CategoryControllerSuite) TestGetAll_StoreError() {
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(nil, 0, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	suite.controller.GetAll(rec, req)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *CategoryControllerSuite) TestQuery() {
	var gotQuery listopt.Query
	// ctx + 3 opts (sort, filter, pagination).
	suite.mockSvc.EXPECT().
		GetAll(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, query listopt.Query) ([]category.Category, int, error) {
			gotQuery = query
			return []category.Category{{Slug: "strategy"}}, 1, nil
		})

	body := `{"pagination":{"page":2,"pageSize":20},"sort":{"field":"name","order":"asc"},"filters":[{"field":"name","op":"like","value":"str"}]}`
	req := httptest.NewRequest("QUERY", "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, req)
	suite.Equal(http.StatusOK, rec.Code)

	suite.Equal("name", gotQuery.Sort.Column)
	suite.Equal("asc", gotQuery.Sort.Order)
	suite.Require().Len(gotQuery.Filters, 1)
	suite.Equal(listopt.Like, gotQuery.Filters[0].Operator)
	suite.Equal(2, gotQuery.Pagination.Page)
	suite.Equal(20, gotQuery.Pagination.PageSize)

	var resp PaginatedResponse[category.Category]
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &resp))
	suite.Equal(1, resp.TotalItems)
	suite.Equal(1, resp.TotalPages)
	suite.Len(resp.Data, 1)
}

func (suite *CategoryControllerSuite) TestQuery_InvalidJSON() {
	req := httptest.NewRequest("QUERY", "/", strings.NewReader(`{`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Query(rec, req)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *CategoryControllerSuite) TestDelete() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "strategy", true).Return(nil)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/?hard=true", nil), "slug", "strategy")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Equal(http.StatusNoContent, rec.Code)
}

func (suite *CategoryControllerSuite) TestDelete_NotFound() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "missing", false).Return(category.ErrNotFound)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/", nil), "slug", "missing")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *CategoryControllerSuite) TestDelete_StillReferenced() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "strategy", false).Return(category.ErrStillReferenced)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/", nil), "slug", "strategy")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Equal(http.StatusConflict, rec.Code)
}

func (suite *CategoryControllerSuite) TestDelete_StoreError() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "strategy", false).Return(assert.AnError)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/", nil), "slug", "strategy")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestCategoryControllerSuite(t *testing.T) {
	suite.Run(t, new(CategoryControllerSuite))
}
