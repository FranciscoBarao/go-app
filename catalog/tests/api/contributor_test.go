package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/store"
	"github.com/stretchr/testify/assert"
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
		Return(contributor.Contributor{}, store.ErrNotFound)

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

func (suite *ContributorSuite) TestPostContributor() {
	suite.base.dbMock.EXPECT().
		GetContributorBySlug(gomock.Any(), "reiner-knizia").
		Return(contributor.Contributor{}, store.ErrNotFound)
	suite.base.dbMock.EXPECT().
		CreateContributor(gomock.Any(), gomock.Any()).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/contributors").
		JSON(`{"name":"Reiner Knizia"}`).
		Header("Content-Type", "application/json").
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(func(res *http.Response, _ *http.Request) error {
			var got contributor.Contributor
			if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
				return err
			}
			suite.Equal("reiner-knizia", got.Slug)
			suite.Equal("Reiner Knizia", got.Name)
			return nil
		}).
		End()
}

func (suite *ContributorSuite) TestPostContributor_DuplicateSlug() {
	suite.base.dbMock.EXPECT().
		GetContributorBySlug(gomock.Any(), "reiner-knizia").
		Return(contributor.Contributor{Slug: "reiner-knizia", Name: "Reiner Knizia"}, nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/contributors").
		JSON(`{"name":"Reiner Knizia"}`).
		Header("Content-Type", "application/json").
		Expect(suite.T()).
		Status(http.StatusConflict).
		Body(`slug "reiner-knizia" already in use; disambiguate the name: entry already exists` + "\n").
		End()
}

func (suite *ContributorSuite) TestPostContributor_EmptyName() {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Post("/api/contributors").
		JSON(`{"name":""}`).
		Header("Content-Type", "application/json").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		End()
}

func (suite *ContributorSuite) TestPatchContributor() {
	suite.base.dbMock.EXPECT().
		GetContributorBySlug(gomock.Any(), "reiner-knizia").
		Return(contributor.Contributor{Slug: "reiner-knizia", Name: "Reiner Knizia"}, nil)
	suite.base.dbMock.EXPECT().
		UpdateContributor(gomock.Any(), gomock.Any()).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Patch("/api/contributors/reiner-knizia").
		JSON(`{"name":"Dr. Knizia"}`).
		Header("Content-Type", "application/json").
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(func(res *http.Response, _ *http.Request) error {
			var got contributor.Contributor
			if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
				return err
			}
			suite.Equal("reiner-knizia", got.Slug)
			suite.Equal("Dr. Knizia", got.Name)
			return nil
		}).
		End()
}

func (suite *ContributorSuite) TestPatchContributor_NotFound() {
	suite.base.dbMock.EXPECT().
		GetContributorBySlug(gomock.Any(), "missing").
		Return(contributor.Contributor{}, store.ErrNotFound)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Patch("/api/contributors/missing").
		JSON(`{"name":"Dr. Knizia"}`).
		Header("Content-Type", "application/json").
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ContributorSuite) TestPatchContributor_EmptyName() {
	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Patch("/api/contributors/reiner-knizia").
		JSON(`{"name":"   "}`).
		Header("Content-Type", "application/json").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		Body("name cannot be empty: validation failed\n").
		End()
}

func (suite *ContributorSuite) TestDeleteContributor() {
	suite.base.dbMock.EXPECT().
		DeleteContributor(gomock.Any(), "reiner-knizia", false).
		Return(nil)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/contributors/reiner-knizia").
		Expect(suite.T()).
		Status(http.StatusNoContent).
		End()
}

func (suite *ContributorSuite) TestDeleteContributor_NotFound() {
	suite.base.dbMock.EXPECT().
		DeleteContributor(gomock.Any(), "missing", false).
		Return(store.ErrNotFound)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/contributors/missing").
		Expect(suite.T()).
		Status(http.StatusNotFound).
		End()
}

func (suite *ContributorSuite) TestDeleteContributor_StillReferenced() {
	suite.base.dbMock.EXPECT().
		DeleteContributor(gomock.Any(), "reiner-knizia", false).
		Return(store.ErrForeignKey)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Delete("/api/contributors/reiner-knizia").
		Expect(suite.T()).
		Status(http.StatusConflict).
		Body(`contributor "reiner-knizia": still referenced` + "\n").
		End()
}

func (suite *ContributorSuite) TestGetContributors_StoreError() {
	suite.base.dbMock.EXPECT().
		GetAllContributors(gomock.Any(), gomock.Any()).
		Return(nil, 0, assert.AnError)

	apitest.New().
		HandlerFunc(suite.base.router.ServeHTTP).
		Get("/api/contributors").
		Expect(suite.T()).
		Status(http.StatusInternalServerError).
		Body("list contributors: internal error\n").
		End()
}

func TestContributorSuite(t *testing.T) {
	suite.Run(t, new(ContributorSuite))
}
