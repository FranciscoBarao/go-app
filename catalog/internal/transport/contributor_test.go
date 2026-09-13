package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ContributorControllerSuite struct {
	suite.Suite
	mockSvc    *MockContributorService
	controller *ContributorController
}

func (suite *ContributorControllerSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockSvc = NewMockContributorService(ctrl)
	suite.controller = NewContributorController(suite.mockSvc)
}

func (suite *ContributorControllerSuite) TestQuery() {
	suite.mockSvc.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return([]contributor.Contributor{{Slug: "reiner-knizia"}}, 1, nil)

	req := httptest.NewRequest("QUERY", "/", strings.NewReader(`{"pagination":{"page":1,"pageSize":10}}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, req)
	suite.Equal(http.StatusOK, rec.Code)

	var resp PaginatedResponse[contributor.Contributor]
	suite.Require().NoError(json.Unmarshal(rec.Body.Bytes(), &resp))
	suite.Equal(1, resp.TotalItems)
	suite.Equal(1, resp.Page)
	suite.Equal(10, resp.PageSize)
	suite.Len(resp.Data, 1)
}

func (suite *ContributorControllerSuite) TestQuery_InvalidFilterOp() {
	req := httptest.NewRequest("QUERY", "/", strings.NewReader(`{"filters":[{"field":"name","op":"bogus","value":"x"}]}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	suite.controller.Query(rec, req)
	suite.Equal(http.StatusUnprocessableEntity, rec.Code)
}

func TestContributorControllerSuite(t *testing.T) {
	suite.Run(t, new(ContributorControllerSuite))
}
