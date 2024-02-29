package tests

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
)

type MechanismSuite struct {
	suite.Suite

	base *Base
}

func (suite *MechanismSuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *MechanismSuite) TestPostMechanism(t *testing.T) {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": "test"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Body(`{"name": "test"}`).
		Status(http.StatusOK).
		End()
}

func (suite *MechanismSuite) TestGetMechanism(t *testing.T) {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/mechanism/test").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusOK).
		Body(`{"name": "test"}`).
		End()
}

func (suite *MechanismSuite) TestDeleteMechanism(t *testing.T) {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/mechanism/test").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusNoContent).
		End()
}

func (suite *MechanismSuite) TestCreateMechanismFailures(t *testing.T) {
	// Several Json Objects on the body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`[{"name":"a"},{"name":"b"}]`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{name:"a"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": 1}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"test": "test"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(``).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Invalid Struct -> NOT maxstringlength(30)
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> NOT alphanum
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": "test.?"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusForbidden).
		End()
}

func (suite *MechanismSuite) TestGetMechanismFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/mechanism/test").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func (suite *MechanismSuite) TestDeleteMechanismFailure(t *testing.T) {
	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/mechanism/test").
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestMechanismSuite(t *testing.T) {
	suite.Run(t, new(MechanismSuite))
}
