package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type MechanismControllerSuite struct {
	suite.Suite
	mockSvc    *MockMechanismService
	controller *MechanismController
}

func (suite *MechanismControllerSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockSvc = NewMockMechanismService(ctrl)
	suite.controller = NewMechanismController(suite.mockSvc)
}

func (suite *MechanismControllerSuite) TestCreate() {
	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mechanism.Mechanism{Slug: "trading", Name: "Trading"}, nil)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"name":"Trading"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *MechanismControllerSuite) TestGet() {
	expected := mechanism.Mechanism{Slug: "trading", Name: "Trading"}
	suite.mockSvc.EXPECT().Get(gomock.Any(), "trading").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", "trading")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()
	suite.controller.Get(rec, req)
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *MechanismControllerSuite) TestQuery() {
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return([]mechanism.Mechanism{{Slug: "trading"}}, 1, nil)

	req := httptest.NewRequest("QUERY", "/", bytes.NewReader([]byte(`{"pagination":{"page":1,"pageSize":10}}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, req)
	suite.Equal(http.StatusOK, rec.Code)

	var resp PaginatedResponse[mechanism.Mechanism]
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &resp))
	suite.Equal(1, resp.TotalItems)
	suite.Equal(1, resp.Page)
	suite.Equal(10, resp.PageSize)
	suite.Len(resp.Data, 1)
}

func TestMechanismControllerSuite(t *testing.T) {
	suite.Run(t, new(MechanismControllerSuite))
}
