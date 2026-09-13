package tests

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/middleware"
)

type CategorySuite struct {
	suite.Suite
	base *Base
}

func (suite *CategorySuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *CategorySuite) TestPostCategory() {
	suite.base.dbMock.EXPECT().
		GetCategoryBySlug(gomock.Any(), "test").
		Return(category.Category{}, middleware.NewError(http.StatusNotFound, "record not found"))
	suite.base.dbMock.EXPECT().
		CreateCategory(gomock.Any(), gomock.Any()).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name":"test"}`).
		Header("Content-Type", "application/json").
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *CategorySuite) TestGetCategory() {
	expected := category.Category{Slug: "test", Name: "test"}
	suite.base.dbMock.EXPECT().
		GetCategoryBySlug(gomock.Any(), "test").
		Return(expected, nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/category/test").
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *CategorySuite) TestDeleteCategory() {
	suite.base.dbMock.EXPECT().
		DeleteCategory(gomock.Any(), "test", false).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/category/test").
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()
}

func (suite *CategorySuite) TestGetCategoryFailure() {
	suite.base.dbMock.EXPECT().
		GetCategoryBySlug(gomock.Any(), "test").
		Return(category.Category{}, middleware.NewError(http.StatusNotFound, "Record not found"))

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/category/test").
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *CategorySuite) TestGetCategories() {
	expected := []category.Category{{Slug: "strategy", Name: "Strategy"}}
	suite.base.dbMock.EXPECT().
		GetAllCategories(gomock.Any(), gomock.Any()).
		Return(expected, 25, nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/category").
		QueryParams(map[string]string{"page": "3", "pageSize": "10", "sort": "name.asc"}).
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(assertEnvelope(suite.T(), 1, 3, 10, 25, 3)).
		End()
}

func (suite *CategorySuite) TestQueryCategories() {
	expected := []category.Category{{Slug: "strategy", Name: "Strategy"}}
	suite.base.dbMock.EXPECT().
		GetAllCategories(gomock.Any(), gomock.Any()).
		Return(expected, len(expected), nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Method("QUERY").
		URL("/api/category").
		Body(`{"pagination":{"page":1,"pageSize":10}}`).
		ContentType("application/json").
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(assertEnvelope(suite.T(), 1, 1, 10, 1, 1)).
		End()
}

func TestCategorySuite(t *testing.T) {
	suite.Run(t, new(CategorySuite))
}
