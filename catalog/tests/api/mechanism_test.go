package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/model"
)

type MechanismSuite struct {
	suite.Suite

	base *Base
}

func (suite *MechanismSuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *MechanismSuite) TestPostMechanism() {
	mechName := "test"
	mech := model.NewMechanism(mechName)
	suite.base.dbMock.EXPECT().
		Create(mech).
		Return(nil)

	mechJson, err := json.Marshal(mech)
	suite.Require().NoError(err)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(mechJson).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Body(string(mechJson)).
		Status(http.StatusOK).
		End()
}

func (suite *MechanismSuite) TestGetMechanism() {
	mechName := "test"
	mech := new(model.Mechanism)
	suite.base.dbMock.EXPECT().
		Read(mech, "", "name = ?", mechName).
		Do(func(mech *model.Mechanism, sort, query, field string) error {
			mech.Name = mechName
			return nil
		}).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/mechanism/"+mechName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusOK).
		Body(`{ "name": "` + mechName + `" }`).
		End()
}

func (suite *MechanismSuite) TestDeleteMechanism() {

	mechName := "test"
	mech := new(model.Mechanism)
	suite.base.dbMock.EXPECT().
		Read(mech, "", "name = ?", mechName).
		Return(nil)

	suite.base.dbMock.EXPECT().
		Delete(new(model.Mechanism)).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/mechanism/"+mechName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()
}

func (suite *MechanismSuite) TestPostMechanismFailures() {
	// Several Json Objects on the body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`[{"name":"a"},{"name":"b"}]`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Malformed Json
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{name:"a"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Unmarshall type error
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": 1}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Unknown Field
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"test": "test"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Empty Body
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(``).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()

	// Invalid Struct -> NOT maxstringlength(30)
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusForbidden).
		End()

	// Invalid Struct -> NOT alphanum
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name": "test.?"}`).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusForbidden).
		End()
}

func (suite *MechanismSuite) TestGetMechanismFailure() {
	mechName := "test"
	mech := new(model.Mechanism)
	suite.base.dbMock.EXPECT().
		Read(mech, "", "name = ?", mechName).
		Return(middleware.NewError(http.StatusNotFound, "Mechanism not found with name: "+mechName))

	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/mechanism/"+mechName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *MechanismSuite) TestDeleteMechanismFailure() {
	mechName := "test"
	mech := new(model.Mechanism)
	suite.base.dbMock.EXPECT().
		Read(mech, "", "name = ?", mechName).
		Return(middleware.NewError(http.StatusNotFound, "Mechanism not found with name: "+mechName))

	// Record not found
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/mechanism/"+mechName).
		Header("Authorization", "Bearer "+suite.base.oauthHeader).
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func TestMechanismSuite(t *testing.T) {
	suite.Run(t, new(MechanismSuite))
}
