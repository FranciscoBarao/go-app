package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/model"
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
	category := model.NewCategory(categoryName)
	suite.base.dbMock.EXPECT().
		Create(category).
		Return(nil)

	categoryJson, err := json.Marshal(category)
	suite.Require().NoError(err)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(categoryJson).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Body(string(categoryJson)).
		Status(http.StatusOK).
		End()
}

func (suite *CategorySuite) TestGetCategory() {
	categoryName := "test"
	category := new(model.Category)
	suite.base.dbMock.EXPECT().
		Read(category, "", "name = ?", categoryName).
		Do(func(category *model.Category, sort, query, field string) error {
			category.Name = categoryName
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
	category := new(model.Category)
	suite.base.dbMock.EXPECT().
		Read(category, "", "name = ?", categoryName).
		Return(nil)

	suite.base.dbMock.EXPECT().
		Delete(new(model.Category)).
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
	category := new(model.Category)
	suite.base.dbMock.EXPECT().
		Read(category, "", "name = ?", categoryName).
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
	category := new(model.Category)
	suite.base.dbMock.EXPECT().
		Read(category, "", "name = ?", categoryName).
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
