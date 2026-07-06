package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type BoardgameControllerSuite struct {
	suite.Suite
	mockSvc    *MockBoardgameService
	controller *BoardgameController
}

func (suite *BoardgameControllerSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockSvc = NewMockBoardgameService(ctrl)
	suite.controller = NewBoardgameController(suite.mockSvc)
}

func (suite *BoardgameControllerSuite) TestCreate() {
	req := boardgame.CreateBoardgameRequest{Name: "Catan", MinPlayers: 2, MaxPlayers: 4}
	expected := boardgame.Boardgame{Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}

	suite.mockSvc.EXPECT().Create(gomock.Any(), gomock.Any(), "").DoAndReturn(
		func(_ context.Context, r *boardgame.CreateBoardgameRequest, _ string) (boardgame.Boardgame, error) {
			suite.Equal(req.Name, r.Name)
			suite.Equal(req.MinPlayers, r.MinPlayers)
			return expected, nil
		},
	)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Create(rec, httpReq)
	suite.Assert().Equal(http.StatusOK, rec.Code)
}

func (suite *BoardgameControllerSuite) TestGet() {
	suite.mockSvc.EXPECT().GetBySlug(gomock.Any(), "catan").Return(boardgame.Boardgame{Slug: "catan"}, nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", "catan")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()

	suite.controller.Get(rec, httpReq)
	suite.Assert().Equal(http.StatusOK, rec.Code)
}

func (suite *BoardgameControllerSuite) TestCreate_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	suite.controller.Create(rec, req)
	suite.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func TestBoardgameControllerSuite(t *testing.T) {
	suite.Run(t, new(BoardgameControllerSuite))
}
