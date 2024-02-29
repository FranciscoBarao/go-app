package tests

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
)

type CategorySuite struct {
	suite.Suite

	base *Base
}

func (suite *CategorySuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *CategorySuite) TestPostCategory(t *testing.T) {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": "test"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Body(`{"name": "test"}`).
		Status(http.StatusOK).
		End()
}

func (suite *CategorySuite) TestGetCategory(t *testing.T) {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/category/test").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		Body(`{"name": "test"}`).
		End()
}

func (suite *CategorySuite) TestDeleteCategory(t *testing.T) {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/category/test").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusNoContent).
		End()
}

func (suite *CategorySuite) TestCreateCategoryFailures(t *testing.T) {
	// Several Json Objects on the body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`[{"name":"a"},{"name":"b"}]`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`{name:"a"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": 1}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`{"test": "test"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(``).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Invalid Struct -> NOT maxstringlength(30)
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> NOT alphanum
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/category").
		JSON(`{"name": "test.?"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()
}

func (suite *CategorySuite) TestGetCategoryFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/category/test").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func (suite *CategorySuite) TestDeleteCategoryFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/category/test").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestCategorySuite(t *testing.T) {
	suite.Run(t, new(CategorySuite))
}
