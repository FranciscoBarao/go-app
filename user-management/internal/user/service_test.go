package user

import (
	"context"
	"net/http"
	"testing"

	"github.com/FranciscoBarao/user-management/internal/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UserServiceSuite struct {
	suite.Suite
	mockDB  *MockDatabase
	service *Service
}

func (s *UserServiceSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.mockDB = NewMockDatabase(ctrl)
	s.service = NewService(s.mockDB)
}

func (s *UserServiceSuite) TestRegister() {
	ctx := context.Background()
	u := &User{Username: "john", Email: "john@test.com"}
	s.mockDB.EXPECT().CreateUser(ctx, u).Return(nil)

	err := s.service.Register(ctx, u, "secret123")
	s.Require().NoError(err)
	s.Assert().NotEmpty(u.Password)
}

func (s *UserServiceSuite) TestRegister_AlreadyRegistered() {
	ctx := context.Background()
	u := &User{Username: "john", Email: "john@test.com"}
	s.mockDB.EXPECT().
		CreateUser(ctx, u).
		Return(middleware.NewError(http.StatusConflict, "entry already registered"))

	err := s.service.Register(ctx, u, "secret123")
	s.Require().Error(err)
	s.Assert().ErrorContains(err, "entry already registered")
}

func (s *UserServiceSuite) TestRegister_InternalError() {
	ctx := context.Background()
	u := &User{Username: "john", Email: "john@test.com"}
	s.mockDB.EXPECT().
		CreateUser(ctx, u).
		Return(assert.AnError)

	err := s.service.Register(ctx, u, "secret123")
	s.Require().Error(err)
	s.Assert().ErrorIs(err, assert.AnError)
}

func (s *UserServiceSuite) TestLogin() {
	ctx := context.Background()
	u := &User{Username: "john", Email: "john@test.com"}
	_ = u.HashPassword("secret123")

	s.mockDB.EXPECT().GetUserByUsername(ctx, "john").Return(*u, nil)

	err := s.service.Login(ctx, "john", "secret123")
	s.Require().NoError(err)
}

func (s *UserServiceSuite) TestLogin_WrongPassword() {
	ctx := context.Background()
	u := &User{Username: "john", Email: "john@test.com"}
	_ = u.HashPassword("secret123")

	s.mockDB.EXPECT().GetUserByUsername(ctx, "john").Return(*u, nil)

	err := s.service.Login(ctx, "john", "wrongpass")
	s.Require().Error(err)

	var mr *middleware.MalformedRequest
	s.ErrorAs(err, &mr)
	s.Assert().Equal(http.StatusUnauthorized, mr.GetStatus())
}

func (s *UserServiceSuite) TestLogin_UserNotFound() {
	ctx := context.Background()
	s.mockDB.EXPECT().
		GetUserByUsername(ctx, "missing").
		Return(User{}, middleware.NewError(http.StatusNotFound, "record not found"))

	err := s.service.Login(ctx, "missing", "pass")
	s.Require().Error(err)
	s.Assert().ErrorContains(err, "record not found")
}

func (s *UserServiceSuite) TestGetAll() {
	ctx := context.Background()
	expected := []User{{Username: "john"}, {Username: "jane"}}
	s.mockDB.EXPECT().GetAllUsers(ctx).Return(expected, nil)

	users, err := s.service.GetAll(ctx)
	s.Require().NoError(err)
	s.Assert().Equal(expected, users)
}

func (s *UserServiceSuite) TestGetAll_Empty() {
	ctx := context.Background()
	s.mockDB.EXPECT().GetAllUsers(ctx).Return([]User{}, nil)

	users, err := s.service.GetAll(ctx)
	s.Require().NoError(err)
	s.Assert().Empty(users)
}

func (s *UserServiceSuite) TestGetAll_InternalError() {
	ctx := context.Background()
	s.mockDB.EXPECT().
		GetAllUsers(ctx).
		Return([]User{}, assert.AnError)

	users, err := s.service.GetAll(ctx)
	s.Require().Error(err)
	s.Assert().Empty(users)
	s.Assert().ErrorIs(err, assert.AnError)
}

func (s *UserServiceSuite) TestDelete() {
	ctx := context.Background()
	s.mockDB.EXPECT().DeleteUser(ctx, "john").Return(nil)

	err := s.service.Delete(ctx, "john")
	s.Require().NoError(err)
}

func (s *UserServiceSuite) TestDelete_NotFound() {
	ctx := context.Background()
	s.mockDB.EXPECT().
		DeleteUser(ctx, "missing").
		Return(middleware.NewError(http.StatusNotFound, "record not found"))

	err := s.service.Delete(ctx, "missing")
	s.Require().Error(err)
	s.Assert().ErrorContains(err, "record not found")
}

func (s *UserServiceSuite) TestDelete_InternalError() {
	ctx := context.Background()
	s.mockDB.EXPECT().
		DeleteUser(ctx, "missing").
		Return(assert.AnError)

	err := s.service.Delete(ctx, "missing")
	s.Require().Error(err)
	s.Assert().ErrorIs(err, assert.AnError)
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}
