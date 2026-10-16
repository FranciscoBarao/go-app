package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ContributorControllerSuite struct {
	suite.Suite
	mockSvc    *MockContributorService
	controller *ContributorController
}

func (suite *ContributorControllerSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockSvc = NewMockContributorService(ctrl)
	suite.controller = NewContributorController(suite.mockSvc)
}

func (suite *ContributorControllerSuite) TestCreate() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(contributor.Contributor{Slug: "reiner-knizia", Name: "Reiner Knizia"}, nil)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Reiner Knizia"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Equal(http.StatusOK, rec.Code)

	var got contributor.Contributor
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &got))
	suite.Equal("reiner-knizia", got.Slug)
}

func (suite *ContributorControllerSuite) TestCreate_EmptyName() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *ContributorControllerSuite) TestCreate_AlreadyExists() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(contributor.Contributor{}, contributor.ErrAlreadyExists)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Reiner Knizia"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Equal(http.StatusConflict, rec.Code)
}

func (suite *ContributorControllerSuite) TestCreate_StoreError() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(contributor.Contributor{}, assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Reiner Knizia"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *ContributorControllerSuite) TestGet() {
	expected := contributor.Contributor{Slug: "reiner-knizia", Name: "Reiner Knizia"}
	suite.mockSvc.EXPECT().Get(gomock.Any(), "reiner-knizia").Return(expected, nil)

	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "slug", "reiner-knizia")
	rec := httptest.NewRecorder()
	suite.controller.Get(rec, req)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *ContributorControllerSuite) TestGet_NotFound() {
	suite.mockSvc.EXPECT().Get(gomock.Any(), "missing").Return(contributor.Contributor{}, contributor.ErrNotFound)

	req := withRouteParam(httptest.NewRequest(http.MethodGet, "/", nil), "slug", "missing")
	rec := httptest.NewRecorder()
	suite.controller.Get(rec, req)
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *ContributorControllerSuite) TestGetAll_MalformedPage() {
	req := httptest.NewRequest(http.MethodGet, "/?page=nope", nil)
	rec := httptest.NewRecorder()
	suite.controller.GetAll(rec, req)
	suite.Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (suite *ContributorControllerSuite) TestGetAll_StoreError() {
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(nil, 0, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	suite.controller.GetAll(rec, req)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *ContributorControllerSuite) TestQuery() {
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return([]contributor.Contributor{{Slug: "reiner-knizia"}}, 1, nil)

	req := httptest.NewRequest("QUERY", "/", strings.NewReader(`{"pagination":{"page":1,"pageSize":10}}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, req)
	suite.Equal(http.StatusOK, rec.Code)

	var resp PaginatedResponse[contributor.Contributor]
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &resp))
	suite.Equal(1, resp.TotalItems)
	suite.Equal(1, resp.Page)
	suite.Equal(10, resp.PageSize)
	suite.Len(resp.Data, 1)
}

func (suite *ContributorControllerSuite) TestQuery_InvalidJSON() {
	req := httptest.NewRequest("QUERY", "/", strings.NewReader(`{`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Query(rec, req)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *ContributorControllerSuite) TestQuery_InvalidFilterOp() {
	req := httptest.NewRequest("QUERY", "/", strings.NewReader(`{"filters":[{"field":"name","op":"bogus","value":"x"}]}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, req)
	suite.Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (suite *ContributorControllerSuite) TestUpdate() {
	name := "Dr. Knizia"
	suite.mockSvc.EXPECT().Update(gomock.Any(), gomock.Any(), "reiner-knizia").Return(contributor.Contributor{Slug: "reiner-knizia", Name: name}, nil)

	req := withRouteParam(httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"name":"Dr. Knizia"}`)), "slug", "reiner-knizia")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Update(rec, req)
	suite.Equal(http.StatusOK, rec.Code)

	var got contributor.Contributor
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &got))
	suite.Equal(name, got.Name)
}

func (suite *ContributorControllerSuite) TestUpdate_EmptyName() {
	req := withRouteParam(httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"name":"   "}`)), "slug", "reiner-knizia")
	req.Header.Set("Content-Type", "application/json")
	suite.mockSvc.EXPECT().Update(gomock.Any(), gomock.Any(), "reiner-knizia").Return(contributor.Contributor{}, contributor.ErrValidation)

	rec := httptest.NewRecorder()
	suite.controller.Update(rec, req)
	suite.Equal(http.StatusBadRequest, rec.Code)
}

func (suite *ContributorControllerSuite) TestUpdate_NotFound() {
	suite.mockSvc.EXPECT().Update(gomock.Any(), gomock.Any(), "missing").Return(contributor.Contributor{}, contributor.ErrNotFound)

	req := withRouteParam(httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"name":"Dr. Knizia"}`)), "slug", "missing")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Update(rec, req)
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *ContributorControllerSuite) TestUpdate_StoreError() {
	suite.mockSvc.EXPECT().Update(gomock.Any(), gomock.Any(), "reiner-knizia").Return(contributor.Contributor{}, assert.AnError)

	req := withRouteParam(httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"name":"Dr. Knizia"}`)), "slug", "reiner-knizia")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Update(rec, req)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *ContributorControllerSuite) TestDelete() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "reiner-knizia", true).Return(nil)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/?hard=true", nil), "slug", "reiner-knizia")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Equal(http.StatusNoContent, rec.Code)
}

func (suite *ContributorControllerSuite) TestDelete_NotFound() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "missing", false).Return(contributor.ErrNotFound)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/", nil), "slug", "missing")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Equal(http.StatusNotFound, rec.Code)
}

func (suite *ContributorControllerSuite) TestDelete_StillReferenced() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "reiner-knizia", false).Return(contributor.ErrStillReferenced)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/", nil), "slug", "reiner-knizia")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Equal(http.StatusConflict, rec.Code)
}

func (suite *ContributorControllerSuite) TestDelete_StoreError() {
	suite.mockSvc.EXPECT().Delete(gomock.Any(), "reiner-knizia", false).Return(assert.AnError)

	req := withRouteParam(httptest.NewRequest(http.MethodDelete, "/", nil), "slug", "reiner-knizia")
	rec := httptest.NewRecorder()
	suite.controller.Delete(rec, req)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func TestContributorControllerSuite(t *testing.T) {
	suite.Run(t, new(ContributorControllerSuite))
}
