package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"

	"github.com/FranciscoBarao/catalog/category"
	"github.com/FranciscoBarao/catalog/middleware"
)

type CategorySuite struct {
	suite.Suite

	base *Base
}

func (suite *CategorySuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *CategorySuite) TestPostCategory() {
	categoryName := "test"
	categoryObj := category.NewCategory(categoryName)
	suite.base.dbMock.EXPECT().
		Create(categoryObj).
		Return(nil)

	categoryJSON, err := json.Marshal(categoryObj)
	suite.Require().NoError(err)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(categoryJSON).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Body(string(categoryJSON)).
		Status(http.StatusOK).
		End()
}

func (suite *CategorySuite) TestGetCategory() {
	categoryName := "test"
	categoryObj := new(category.Category)
	suite.base.dbMock.EXPECT().
		Read(categoryObj, "", "name = ?", categoryName).
		Do(func(categoryObj *category.Category, sort, query, field string) error {
			categoryObj.Name = categoryName
			return nil
		}).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/category/"+categoryName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(`{"name": "` + categoryName + `"}`).
		End()
}

func (suite *CategorySuite) TestDeleteCategory() {
	categoryName := "test"
	categoryObj := new(category.Category)
	suite.base.dbMock.EXPECT().
		Read(categoryObj, "", "name = ?", categoryName).
		Return(nil)

	suite.base.dbMock.EXPECT().
		Delete(new(category.Category)).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/category/"+categoryName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()
}

func (suite *CategorySuite) TestPostCategoryFailures() {
	// Several Json Objects on the body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`[{"name":"a"},{"name":"b"}]`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`{name:"a"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": 1}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`{"test": "test"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(``).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Invalid Struct -> NOT maxstringlength(30)
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> NOT alphanum
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": "test.?"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusForbidden).
		End()
}

func (suite *CategorySuite) TestGetCategoryFailure() {
	categoryName := "test"
	categoryObj := new(category.Category)
	suite.base.dbMock.EXPECT().
		Read(categoryObj, "", "name = ?", categoryName).
		Return(middleware.NewError(http.StatusNotFound, "Category not found with name: "+categoryName))

	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/category/"+categoryName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *CategorySuite) TestDeleteCategoryFailure() {
	categoryName := "test"
	categoryObj := new(category.Category)
	suite.base.dbMock.EXPECT().
		Read(categoryObj, "", "name = ?", categoryName).
		Return(middleware.NewError(http.StatusNotFound, "Category not found with name: "+categoryName))

	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/category/"+categoryName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func TestCategorySuite(t *testing.T) {
	suite.Run(t, new(CategorySuite))
}
