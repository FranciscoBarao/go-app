package transport

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/category"
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

func TestCategoryControllerSuite(t *testing.T) {
	suite.Run(t, new(CategoryControllerSuite))
}
