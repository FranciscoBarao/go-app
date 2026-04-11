package tests

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/tag"
)

type TagSuite struct {
	suite.Suite

	base *Base
}

func (suite *TagSuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *TagSuite) TestPost() {
	tagName := "test"
	tagObj := tag.NewTag(tagName)
	suite.base.dbMock.EXPECT().
		CreateTag(gomock.Any(), tagObj).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name":"test"}`).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *TagSuite) TestGet() {
	tagName := "test"
	expected := tag.Tag{Name: tagName}
	suite.base.dbMock.EXPECT().
		GetTag(gomock.Any(), tagName).
		Return(expected, nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/tag/"+tagName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *TagSuite) TestDelete() {
	tagName := "test"
	suite.base.dbMock.EXPECT().
		DeleteTag(gomock.Any(), tagName).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/tag/"+tagName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()
}

func (suite *TagSuite) TestPostFailures() {
	// Several Json Objects on the body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/tag").
		JSON(`[{"name":"a"},{"name":"b"}]`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/tag").
		JSON(`{name:"a"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name": 1}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"test": "test"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/tag").
		JSON(``).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Invalid Struct -> NOT maxstringlength(30)
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> NOT alphanum
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name": "test.?"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusForbidden).
		End()
}

func (suite *TagSuite) TestGetFailure() {
	tagName := "test"
	suite.base.dbMock.EXPECT().
		GetTag(gomock.Any(), tagName).
		Return(tag.Tag{}, middleware.NewError(http.StatusNotFound, "Record not found"))

	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/tag/"+tagName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *TagSuite) TestDeleteFailure() {
	tagName := "test"
	suite.base.dbMock.EXPECT().
		DeleteTag(gomock.Any(), tagName).
		Return(middleware.NewError(http.StatusNotFound, "Record not found"))

	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/tag/"+tagName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func TestTagSuite(t *testing.T) {
	suite.Run(t, new(TagSuite))
}
