package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/tag"
	"github.com/go-chi/chi/v5"
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
	suite.mockSvc.EXPECT().Create(gomock.Any()).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Strategy"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result tag.Tag
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal("Strategy", result.Name)
}

func (suite *TagControllerSuite) TestGetAll() {
	expected := []tag.Tag{{Name: "Strategy"}, {Name: "Family"}}
	suite.mockSvc.EXPECT().GetAll("").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	suite.controller.GetAll(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result []tag.Tag
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Len(result, 2)
}

func (suite *TagControllerSuite) TestGet() {
	suite.mockSvc.EXPECT().Get("Strategy").Return(tag.Tag{Name: "Strategy"}, nil)

	req := reqWithParam(httptest.NewRequest(http.MethodGet, "/", nil), "name", "Strategy")
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	var result tag.Tag
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	suite.Equal("Strategy", result.Name)
}

func (suite *TagControllerSuite) TestDelete() {
	suite.mockSvc.EXPECT().Delete("Strategy").Return(nil)

	req := reqWithParam(httptest.NewRequest(http.MethodDelete, "/", nil), "name", "Strategy")
	rec := httptest.NewRecorder()

	suite.controller.Delete(rec, req)

	suite.Equal(http.StatusNoContent, rec.Code)
}

func (suite *TagControllerSuite) TestCreateInvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *TagControllerSuite) TestServiceError() {
	suite.mockSvc.EXPECT().Get("Unknown").Return(tag.Tag{}, middleware.NewError(http.StatusNotFound, "not found"))

	req := reqWithParam(httptest.NewRequest(http.MethodGet, "/", nil), "name", "Unknown")
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func TestTagControllerSuite(t *testing.T) {
	suite.Run(t, new(TagControllerSuite))
}
