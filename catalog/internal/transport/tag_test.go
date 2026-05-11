package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/tag"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type TagControllerSuite struct {
	suite.Suite
	mockSvc    *MockTagService
	controller *TagController
}

func (suite *TagControllerSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockSvc = NewMockTagService(ctrl)
	suite.controller = NewTagController(suite.mockSvc)
}

// reqWithParam creates a request with a chi URL param injected into context.
func reqWithParam(r *http.Request, key, val string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, val)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func (suite *TagControllerSuite) TestCreate() {
	expectedTag := &tag.Tag{Name: "Strategy"}
	suite.mockSvc.EXPECT().Create(gomock.Any(), expectedTag).Return(nil)

	tagBytes, err := json.Marshal(expectedTag)
	suite.Require().NoError(err)
	body := bytes.NewReader(tagBytes)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result tag.Tag
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal(expectedTag.Name, result.Name)
}

func (suite *TagControllerSuite) TestCreate_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *TagControllerSuite) TestCreate_InvalidStruct() {
	const maxChars = 30
	longName := strings.Repeat("a", maxChars+1)
	expectedTag := &tag.Tag{
		Name: longName, // invalid name length
	}

	tagBytes, err := json.Marshal(expectedTag)
	suite.Require().NoError(err)
	body := bytes.NewReader(tagBytes)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *TagControllerSuite) TestCreate_InternalError() {
	expectedTag := &tag.Tag{Name: "Strategy"}

	suite.mockSvc.EXPECT().
		Create(gomock.Any(), expectedTag).
		Return(assert.AnError)

	tagBytes, err := json.Marshal(expectedTag)
	suite.Require().NoError(err)
	body := bytes.NewReader(tagBytes)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TagControllerSuite) TestGetAll() {
	expected := []tag.Tag{{Name: "Strategy"}, {Name: "Family"}}
	suite.mockSvc.EXPECT().GetAll(gomock.Any()).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result []tag.Tag
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Len(result, 2)
}

func (suite *TagControllerSuite) TestGetAll_Empty() {
	expected := []tag.Tag{}
	suite.mockSvc.EXPECT().GetAll(gomock.Any()).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result []tag.Tag
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Empty(result)
}

func (suite *TagControllerSuite) TestGetAll_InternalError() {
	suite.mockSvc.EXPECT().
		GetAll(gomock.Any()).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TagControllerSuite) TestGet() {
	expectedTag := tag.Tag{Name: "Strategy"}
	suite.mockSvc.EXPECT().Get(gomock.Any(), expectedTag.Name).Return(expectedTag, nil)

	req := reqWithParam(
		httptest.NewRequest(http.MethodGet, "/", nil),
		"name", expectedTag.Name,
	)
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result tag.Tag
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal(expectedTag.Name, result.Name)
}

func (suite *TagControllerSuite) TestGet_NotFound() {
	name := "not found"
	suite.mockSvc.EXPECT().
		Get(gomock.Any(), name).
		Return(tag.Tag{}, middleware.NewError(http.StatusNotFound, "not found"))

	req := reqWithParam(
		httptest.NewRequest(http.MethodGet, "/", nil),
		"name", name,
	)
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *TagControllerSuite) TestGet_InternalError() {
	name := "not found"
	suite.mockSvc.EXPECT().
		Get(gomock.Any(), name).
		Return(tag.Tag{}, assert.AnError)

	req := reqWithParam(
		httptest.NewRequest(http.MethodGet, "/", nil),
		"name", name,
	)
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *TagControllerSuite) TestDelete() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "Strategy").Return(nil)

	req := reqWithParam(httptest.NewRequest(http.MethodDelete, "/", nil), "name", "Strategy")
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusNoContent, rec.Code)
}

func (suite *TagControllerSuite) TestDelete_NotFound() {
	name := "not found"
	suite.mockSvc.EXPECT().
		Delete(gomock.Any(), name).
		Return(middleware.NewError(http.StatusNotFound, "not found"))

	req := reqWithParam(
		httptest.NewRequest(http.MethodDelete, "/", nil),
		"name", name,
	)
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *TagControllerSuite) TestDelete_InternalError() {
	name := "Strategy"
	suite.mockSvc.EXPECT().
		Delete(gomock.Any(), name).
		Return(assert.AnError)

	req := reqWithParam(
		httptest.NewRequest(http.MethodDelete, "/", nil),
		"name", name,
	)
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestTagControllerSuite(t *testing.T) {
	suite.Run(t, new(TagControllerSuite))
}
