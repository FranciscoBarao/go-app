package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"

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
		Create(tagObj).
		Return(nil)

	tagJSON, err := json.Marshal(tagObj)
	suite.Require().NoError(err)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/tag").
		JSON(tagJSON).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Body(string(tagJSON)).
		Status(http.StatusOK).
		End()
}

func (suite *TagSuite) TestGet() {
	tagName := "test"
	tagObj := new(tag.Tag)
	suite.base.dbMock.EXPECT().
		Read(tagObj, "", "name = ?", tagName).
		Do(func(t *tag.Tag, sort, query, field string) error {
			t.Name = tagName
			return nil
		}).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/tag/"+tagName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(`{ "name": "` + tagName + `" }`).
		End()
}

func (suite *TagSuite) TestDelete() {
	tagName := "test"
	tagObj := new(tag.Tag)
	suite.base.dbMock.EXPECT().
		Read(tagObj, "", "name = ?", tagName).
		Return(nil)

	suite.base.dbMock.EXPECT().
		Delete(new(tag.Tag)).
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
	tagObj := new(tag.Tag)
	suite.base.dbMock.EXPECT().
		Read(tagObj, "", "name = ?", tagName).
		Return(middleware.NewError(http.StatusNotFound, "Tag not found with name: "+tagName))

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
	tagObj := new(tag.Tag)
	suite.base.dbMock.EXPECT().
		Read(tagObj, "", "name = ?", tagName).
		Return(middleware.NewError(http.StatusNotFound, "Tag not found with name: "+tagName))

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
