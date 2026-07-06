package tests

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/middleware"
)

type MechanismSuite struct {
	suite.Suite
	base *Base
}

func (suite *MechanismSuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *MechanismSuite) TestPostMechanism() {
	suite.base.dbMock.EXPECT().
		GetMechanismBySlug(gomock.Any(), "test").
		Return(mechanism.Mechanism{}, middleware.NewError(http.StatusNotFound, "record not found"))
	suite.base.dbMock.EXPECT().
		CreateMechanism(gomock.Any(), gomock.Any()).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name":"test"}`).
		Header("Content-Type", "application/json").
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *MechanismSuite) TestGetMechanism() {
	expected := mechanism.Mechanism{Slug: "test", Name: "test"}
	suite.base.dbMock.EXPECT().
		GetMechanismBySlug(gomock.Any(), "test").
		Return(expected, nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/mechanism/test").
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *MechanismSuite) TestDeleteMechanism() {
	suite.base.dbMock.EXPECT().
		DeleteMechanism(gomock.Any(), "test", false).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/mechanism/test").
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()
}

func (suite *MechanismSuite) TestGetMechanismFailure() {
	suite.base.dbMock.EXPECT().
		GetMechanismBySlug(gomock.Any(), "test").
		Return(mechanism.Mechanism{}, middleware.NewError(http.StatusNotFound, "Record not found"))

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/mechanism/test").
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func TestMechanismSuite(t *testing.T) {
	suite.Run(t, new(MechanismSuite))
}
