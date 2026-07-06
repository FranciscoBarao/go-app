package transport

import (
	"bytes"
	"context"
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

func TestMechanismControllerSuite(t *testing.T) {
	suite.Run(t, new(MechanismControllerSuite))
}
