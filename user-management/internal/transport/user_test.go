package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/user-management/internal/middleware"
	"github.com/FranciscoBarao/user-management/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UserControllerSuite struct {
	suite.Suite
	mockSvc    *MockUserService
	controller *UserController
}

func (s *UserControllerSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.mockSvc = NewMockUserService(ctrl)
	s.controller = NewUserController(s.mockSvc)
}

func reqWithParam(r *http.Request, key, val string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, val)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// --- Register tests ---

func (s *UserControllerSuite) TestRegister() {
	req := user.RegisterRequest{Username: "john", Email: "john@test.com", Password: "secret123"}
	expectedUser := user.NewUser(&req)

	s.mockSvc.EXPECT().Register(gomock.Any(), expectedUser, "secret123").Return(nil)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.controller.Register(rec, httpReq)

	s.Assert().Equal(http.StatusCreated, rec.Code)
	var result user.User
	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	s.Assert().Equal("john", result.Username)
	s.Assert().Equal("john@test.com", result.Email)
}

func (s *UserControllerSuite) TestRegister_InvalidJSON() {
	httpReq := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.controller.Register(rec, httpReq)

	s.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (s *UserControllerSuite) TestRegister_ValidationError() {
	req := user.RegisterRequest{Username: "", Email: "bad", Password: "short"}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.controller.Register(rec, httpReq)

	s.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (s *UserControllerSuite) TestRegister_ServiceError() {
	req := user.RegisterRequest{Username: "john", Email: "john@test.com", Password: "secret123"}
	expectedUser := user.NewUser(&req)

	s.mockSvc.EXPECT().Register(gomock.Any(), expectedUser, "secret123").Return(assert.AnError)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.controller.Register(rec, httpReq)

	s.Assert().Equal(http.StatusInternalServerError, rec.Code)
}

func (s *UserControllerSuite) TestRegister_Conflict() {
	req := user.RegisterRequest{Username: "john", Email: "john@test.com", Password: "secret123"}
	expectedUser := user.NewUser(&req)

	s.mockSvc.EXPECT().Register(gomock.Any(), expectedUser, "secret123").Return(middleware.NewError(http.StatusConflict, "entry already registered"))

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.controller.Register(rec, httpReq)

	s.Assert().Equal(http.StatusConflict, rec.Code)
}

// --- GetAll tests ---

func (s *UserControllerSuite) TestGetAll() {
	expected := []user.User{{Username: "john"}, {Username: "jane"}}
	s.mockSvc.EXPECT().GetAll(gomock.Any()).Return(expected, nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	s.controller.GetAll(rec, httpReq)

	s.Assert().Equal(http.StatusOK, rec.Code)
	var result []user.User
	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	s.Assert().Len(result, 2)
}

func (s *UserControllerSuite) TestGetAll_Empty() {
	s.mockSvc.EXPECT().GetAll(gomock.Any()).Return([]user.User{}, nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	s.controller.GetAll(rec, httpReq)

	s.Assert().Equal(http.StatusOK, rec.Code)
	var result []user.User
	s.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &result))
	s.Assert().Empty(result)
}

func (s *UserControllerSuite) TestGetAll_InternalError() {
	s.mockSvc.EXPECT().GetAll(gomock.Any()).Return(nil, assert.AnError)

	httpReq := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	s.controller.GetAll(rec, httpReq)

	s.Assert().Equal(http.StatusInternalServerError, rec.Code)
}

// --- Delete tests ---

func (s *UserControllerSuite) TestDelete() {
	s.mockSvc.EXPECT().Delete(gomock.Any(), "john").Return(nil)

	httpReq := reqWithParam(httptest.NewRequest(http.MethodDelete, "/", nil), "username", "john")
	rec := httptest.NewRecorder()

	s.controller.Delete(rec, httpReq)

	s.Assert().Equal(http.StatusNoContent, rec.Code)
}

func (s *UserControllerSuite) TestDelete_InvalidUsername() {
	// Username with spaces/special chars fails alphanum validation
	httpReq := reqWithParam(httptest.NewRequest(http.MethodDelete, "/", nil), "username", "bad user!")
	rec := httptest.NewRecorder()

	s.controller.Delete(rec, httpReq)

	s.Assert().Equal(http.StatusBadRequest, rec.Code)
}

func (s *UserControllerSuite) TestDelete_NotFound() {
	s.mockSvc.EXPECT().Delete(gomock.Any(), "missing").Return(middleware.NewError(http.StatusNotFound, "record not found"))

	httpReq := reqWithParam(httptest.NewRequest(http.MethodDelete, "/", nil), "username", "missing")
	rec := httptest.NewRecorder()

	s.controller.Delete(rec, httpReq)

	s.Assert().Equal(http.StatusNotFound, rec.Code)
}

func (s *UserControllerSuite) TestDelete_InternalError() {
	s.mockSvc.EXPECT().Delete(gomock.Any(), "john").Return(assert.AnError)

	httpReq := reqWithParam(httptest.NewRequest(http.MethodDelete, "/", nil), "username", "john")
	rec := httptest.NewRecorder()

	s.controller.Delete(rec, httpReq)

	s.Assert().Equal(http.StatusInternalServerError, rec.Code)
}

func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, new(UserControllerSuite))
}
