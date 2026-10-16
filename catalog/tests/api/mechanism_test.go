package tests

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/store"
	"github.com/stretchr/testify/assert"
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
		Return(mechanism.Mechanism{}, store.ErrNotFound)
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
		Return(mechanism.Mechanism{}, store.ErrNotFound)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/mechanism/test").
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *MechanismSuite) TestGetMechanisms() {
	expected := []mechanism.Mechanism{{Slug: "trading", Name: "Trading"}}
	suite.base.dbMock.EXPECT().
		GetAllMechanisms(gomock.Any(), gomock.Any()).
		Return(expected, len(expected), nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/mechanism").
		QueryParams(map[string]string{"filter": "name.like.trad"}).
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(assertEnvelope(suite.T(), 1, 1, 10, 1, 1)).
		End()
}

func (suite *MechanismSuite) TestQueryMechanisms() {
	expected := []mechanism.Mechanism{{Slug: "trading", Name: "Trading"}}
	suite.base.dbMock.EXPECT().
		GetAllMechanisms(gomock.Any(), gomock.Any()).
		Return(expected, len(expected), nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Method("QUERY").
		URL("/api/mechanism").
		Body(`{"pagination":{"page":1,"pageSize":10}}`).
		ContentType("application/json").
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(assertEnvelope(suite.T(), 1, 1, 10, 1, 1)).
		End()
}

func (suite *MechanismSuite) TestPostMechanism_DuplicateSlug() {
	suite.base.dbMock.EXPECT().
		GetMechanismBySlug(gomock.Any(), "trading").
		Return(mechanism.Mechanism{Slug: "trading", Name: "Trading"}, nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name":"Trading"}`).
		Header("Content-Type", "application/json").
		Expect(suite.T()).
		Status(http.StatusConflict).
		Body(`slug "trading" already in use; disambiguate the name: entry already exists` + "\n").
		End()
}

func (suite *MechanismSuite) TestPostMechanism_EmptyName() {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/mechanism").
		JSON(`{"name":""}`).
		Header("Content-Type", "application/json").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *MechanismSuite) TestDeleteMechanism_NotFound() {
	suite.base.dbMock.EXPECT().
		DeleteMechanism(gomock.Any(), "missing", false).
		Return(store.ErrNotFound)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/mechanism/missing").
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *MechanismSuite) TestDeleteMechanism_StillReferenced() {
	suite.base.dbMock.EXPECT().
		DeleteMechanism(gomock.Any(), "trading", false).
		Return(store.ErrForeignKey)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/mechanism/trading").
		Expect(suite.T()).
		Status(http.StatusConflict).
		Body(`mechanism "trading": still referenced` + "\n").
		End()
}

func (suite *MechanismSuite) TestGetMechanisms_MalformedPage() {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/mechanism").
		QueryParams(map[string]string{"page": "nope"}).
		Expect(suite.T()).
		Status(http.StatusUnprocessableEntity).
		End()
}

func (suite *MechanismSuite) TestGetMechanisms_StoreError() {
	suite.base.dbMock.EXPECT().
		GetAllMechanisms(gomock.Any(), gomock.Any()).
		Return(nil, 0, assert.AnError)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/mechanism").
		Expect(suite.T()).
		Status(http.StatusInternalServerError).
		Body("list mechanisms: internal error\n").
		End()
}

func TestMechanismSuite(t *testing.T) {
	suite.Run(t, new(MechanismSuite))
}
