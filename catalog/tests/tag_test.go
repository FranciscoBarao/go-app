package tests

import (
	"catalog/middleware"
	"catalog/model"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
)

type TagSuite struct {
	suite.Suite

	base *Base
}

func (suite *TagSuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *TagSuite) TestCreateTagSuccess() {
	tagName := "test"
	tag := model.NewTag(tagName)
	suite.base.dbMock.EXPECT().
		Create(tag).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/tag").
		JSON(`{"name": "`+tagName+`"}`).
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Body(`{"name": "` + tagName + `"}`).
		Status(http.StatusOK).
		End()
}

func (suite *TagSuite) TestGetTagSuccess() {
	tagName := "test"
	tag := new(model.Tag)
	suite.base.dbMock.EXPECT().
		Read(tag, "", "name = ?", tagName).
		Do(func(tag *model.Tag, sort, query, field string) error {
			tag.Name = tagName
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

func (suite *TagSuite) TestGetAllTagSuccess() {
	tags := new([]model.Tag)
	suite.base.dbMock.EXPECT().
		Read(tags, "", "", "").
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/tag").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *TagSuite) TestDeleteTagSuccess() {
	tagName := "test"
	tag := new(model.Tag)
	suite.base.dbMock.EXPECT().
		Read(tag, "", "name = ?", tagName).
		Return(nil)

	suite.base.dbMock.EXPECT().
		Delete(new(model.Tag)).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/tag/"+tagName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()
}

func (suite *TagSuite) TestCreateTagFailures() {
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

func (suite *TagSuite) TestGetTagFailure() {
	tagName := "test"
	tag := new(model.Tag)
	suite.base.dbMock.EXPECT().
		Read(tag, "", "name = ?", tagName).
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

func (suite *TagSuite) TestDeleteTagFailure() {
	tagName := "test"
	tag := new(model.Tag)
	suite.base.dbMock.EXPECT().
		Read(tag, "", "name = ?", tagName).
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
