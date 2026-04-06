package category

import (
	"testing"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// CategoryServiceSuite tests the category Service in isolation using a MockDatabase.
type CategoryServiceSuite struct {
	suite.Suite
	mockDB  *MockDatabase
	service *Service
}

func (suite *CategoryServiceSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockDB = NewMockDatabase(ctrl)
	suite.service = NewService(suite.mockDB)
}

func (suite *CategoryServiceSuite) TearDownTest() {}

func (suite *CategoryServiceSuite) TestCreate() {
	cat := NewCategory("strategy")
	suite.mockDB.EXPECT().Create(cat).Return(nil)

	err := suite.service.Create(cat)
	suite.Assert().NoError(err)
}

func (suite *CategoryServiceSuite) TestGetAll() {
	expected := []Category{{Name: "strategy"}, {Name: "cooperative"}}
	suite.mockDB.EXPECT().Read(gomock.Any(), "name", "", "").DoAndReturn(
		func(dest interface{}, sort, search, identifier string) error {
			ptr := dest.(*[]Category)
			*ptr = expected
			return nil
		},
	)

	categories, err := suite.service.GetAll("name")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, categories)
}

func (suite *CategoryServiceSuite) TestGet() {
	expected := Category{Name: "strategy"}
	suite.mockDB.EXPECT().Read(gomock.Any(), "", "name = ?", "strategy").DoAndReturn(
		func(dest interface{}, sort, search, identifier string) error {
			ptr := dest.(*Category)
			*ptr = expected
			return nil
		},
	)

	cat, err := suite.service.Get("strategy")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, cat)
}

func (suite *CategoryServiceSuite) TestGetNotFound() {
	suite.mockDB.EXPECT().Read(gomock.Any(), "", "name = ?", "nonexistent").Return(
		middleware.NewError(404, "not found"),
	)

	_, err := suite.service.Get("nonexistent")
	suite.Assert().Error(err)
	suite.Assert().Equal("Category not found with name: nonexistent", err.Error())
}

func (suite *CategoryServiceSuite) TestDelete() {
	expected := Category{Name: "strategy"}
	suite.mockDB.EXPECT().Read(gomock.Any(), "", "name = ?", "strategy").DoAndReturn(
		func(dest interface{}, sort, search, identifier string) error {
			ptr := dest.(*Category)
			*ptr = expected
			return nil
		},
	)
	suite.mockDB.EXPECT().Delete(&expected).Return(nil)

	err := suite.service.Delete("strategy")
	suite.Assert().NoError(err)
}

func (suite *CategoryServiceSuite) TestDeleteNotFound() {
	suite.mockDB.EXPECT().Read(gomock.Any(), "", "name = ?", "missing").Return(
		middleware.NewError(404, "not found"),
	)
	// db.Delete should NOT be called

	err := suite.service.Delete("missing")
	suite.Assert().Error(err)
	suite.Assert().Equal("Category not found with name: missing", err.Error())
}

func TestCategoryServiceSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceSuite))
}
