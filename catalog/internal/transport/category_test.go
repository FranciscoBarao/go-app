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
	"github.com/go-chi/chi/v5"
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

func (suite *CategoryControllerSuite) TestGet() {
	expected := category.Category{Slug: "strategy", Name: "Strategy"}
	suite.mockSvc.EXPECT().Get(gomock.Any(), "strategy").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", "strategy")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	suite.controller.Get(rec, req)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *CategoryControllerSuite) TestQuery() {
	var gotParams listopt.Params
	// ctx + 3 opts (sort, filter, pagination).
	suite.mockSvc.EXPECT().
		GetAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, opts ...listopt.Option) ([]category.Category, int, error) {
			gotParams = listopt.Apply(opts...)
			return []category.Category{{Slug: "strategy"}}, 1, nil
		})

	body := `{"pagination":{"page":2,"pageSize":20},"sort":{"field":"name","order":"asc"},"filters":[{"field":"name","op":"like","value":"str"}]}`
	req := httptest.NewRequest("QUERY", "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, req)
	suite.Equal(http.StatusOK, rec.Code)

	suite.Equal("name", gotParams.Sort.Column)
	suite.Equal("asc", gotParams.Sort.Order)
	suite.Require().Len(gotParams.Filters, 1)
	suite.Equal(listopt.OpLike, gotParams.Filters[0].Op)
	suite.Equal(2, gotParams.Pagination.Page)
	suite.Equal(20, gotParams.Pagination.PageSize)

	var resp PaginatedResponse[category.Category]
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &resp))
	suite.Equal(1, resp.TotalItems)
	suite.Equal(1, resp.TotalPages)
	suite.Len(resp.Data, 1)
}

func TestCategoryControllerSuite(t *testing.T) {
	suite.Run(t, new(CategoryControllerSuite))
}
