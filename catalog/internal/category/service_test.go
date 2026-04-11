package category

import (
	"context"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// CategoryServiceSuite tests the category Service in isolation using a MockDatabase.
type CategoryServiceSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	mockDB  *MockDatabase
	service *Service
}

func (suite *CategoryServiceSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockDB = NewMockDatabase(suite.ctrl)
	suite.service = NewService(suite.mockDB)
}

func (suite *CategoryServiceSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *CategoryServiceSuite) TestCreate() {
	cat := NewCategory("strategy")
	suite.mockDB.EXPECT().CreateCategory(gomock.Any(), cat).Return(nil)

	err := suite.service.Create(context.Background(), cat)
	suite.Assert().NoError(err)
}

func (suite *CategoryServiceSuite) TestGetAll() {
	expected := []Category{{Name: "strategy"}, {Name: "cooperative"}}
	suite.mockDB.EXPECT().GetAllCategories(gomock.Any(), "name").Return(expected, nil)

	categories, err := suite.service.GetAll(context.Background(), "name")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, categories)
}

func (suite *CategoryServiceSuite) TestGet() {
	expected := Category{Name: "strategy"}
	suite.mockDB.EXPECT().GetCategory(gomock.Any(), "strategy").Return(expected, nil)

	cat, err := suite.service.Get(context.Background(), "strategy")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, cat)
}

func (suite *CategoryServiceSuite) TestGetNotFound() {
	suite.mockDB.EXPECT().GetCategory(gomock.Any(), "nonexistent").Return(
		Category{}, middleware.NewError(404, "Record not found"),
	)

	_, err := suite.service.Get(context.Background(), "nonexistent")
	suite.Assert().Error(err)
	suite.Assert().Equal("Record not found", err.Error())
}

func (suite *CategoryServiceSuite) TestDelete() {
	suite.mockDB.EXPECT().DeleteCategory(gomock.Any(), "strategy").Return(nil)

	err := suite.service.Delete(context.Background(), "strategy")
	suite.Assert().NoError(err)
}

func (suite *CategoryServiceSuite) TestDeleteNotFound() {
	suite.mockDB.EXPECT().DeleteCategory(gomock.Any(), "missing").Return(
		middleware.NewError(404, "Record not found"),
	)

	err := suite.service.Delete(context.Background(), "missing")
	suite.Assert().Error(err)
	suite.Assert().Equal("Record not found", err.Error())
}

func TestCategoryServiceSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceSuite))
}
