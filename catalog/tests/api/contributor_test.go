package tests

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/middleware"
)

type ContributorSuite struct {
	suite.Suite
	base *Base
}

func (suite *ContributorSuite) SetupSuite() {
	suite.base = NewBase(suite.T())
}

func (suite *ContributorSuite) TestGetContributor() {
	expected := contributor.Contributor{Slug: "reiner-knizia", Name: "Reiner Knizia"}
	suite.base.dbMock.EXPECT().
		GetContributorBySlug(gomock.Any(), "reiner-knizia").
		Return(expected, nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/contributors/reiner-knizia").
		Expect(suite.T()).
		Status(http.StatusOK).
		End()
}

func (suite *ContributorSuite) TestGetContributorFailure() {
	suite.base.dbMock.EXPECT().
		GetContributorBySlug(gomock.Any(), "missing").
		Return(contributor.Contributor{}, middleware.NewError(http.StatusNotFound, "Record not found"))

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/contributors/missing").
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ContributorSuite) TestGetContributors() {
	expected := []contributor.Contributor{{Slug: "reiner-knizia", Name: "Reiner Knizia"}}
	suite.base.dbMock.EXPECT().
		GetAllContributors(gomock.Any(), gomock.Any()).
		Return(expected, len(expected), nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/contributors").
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(assertEnvelope(suite.T(), 1, 1, 10, 1, 1)).
		End()
}

func (suite *ContributorSuite) TestQueryContributors() {
	expected := []contributor.Contributor{{Slug: "reiner-knizia", Name: "Reiner Knizia"}}
	suite.base.dbMock.EXPECT().
		GetAllContributors(gomock.Any(), gomock.Any()).
		Return(expected, len(expected), nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Method("QUERY").
		URL("/api/contributors").
		Body(`{"pagination":{"page":1,"pageSize":10}}`).
		ContentType("application/json").
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(assertEnvelope(suite.T(), 1, 1, 10, 1, 1)).
		End()
}

func TestContributorSuite(t *testing.T) {
	suite.Run(t, new(ContributorSuite))
}
