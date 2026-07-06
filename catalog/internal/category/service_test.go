package category

import (
	"context"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

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
	suite.mockDB.EXPECT().GetCategoryBySlug(gomock.Any(), "strategy").Return(Category{}, middleware.NewError(404, "record not found"))
	suite.mockDB.EXPECT().CreateCategory(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, c *Category) error {
			suite.Equal("strategy", c.Slug)
			suite.Equal("Strategy", c.Name)
			return nil
		},
	)

	cat, err := suite.service.Create(context.Background(), &CreateCategoryRequest{Name: "Strategy"})
	suite.Assert().NoError(err)
	suite.Assert().Equal("strategy", cat.Slug)
}

func (suite *CategoryServiceSuite) TestCreate_DuplicateSlug() {
	suite.mockDB.EXPECT().GetCategoryBySlug(gomock.Any(), "strategy").Return(Category{Slug: "strategy", Name: "Strategy"}, nil)

	_, err := suite.service.Create(context.Background(), &CreateCategoryRequest{Name: "Strategy"})
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().ErrorAs(err, &mr)
	suite.Assert().Equal(409, mr.GetStatus())
}

func (suite *CategoryServiceSuite) TestCreate_EmptySlug() {
	_, err := suite.service.Create(context.Background(), &CreateCategoryRequest{Name: "!!!"})
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().ErrorAs(err, &mr)
	suite.Assert().Equal(400, mr.GetStatus())
}

func (suite *CategoryServiceSuite) TestGetAll() {
	expected := []Category{{Slug: "strategy", Name: "Strategy"}}
	suite.mockDB.EXPECT().GetAllCategories(gomock.Any(), listopt.Params{Sort: listopt.Sort{Column: "name", Order: "asc"}}).Return(expected, nil)

	categories, err := suite.service.GetAll(context.Background(), listopt.WithSort("name", "asc"))
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, categories)
}

func (suite *CategoryServiceSuite) TestGet() {
	expected := Category{Slug: "strategy", Name: "Strategy"}
	suite.mockDB.EXPECT().GetCategoryBySlug(gomock.Any(), "strategy").Return(expected, nil)

	cat, err := suite.service.Get(context.Background(), "strategy")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, cat)
}

func (suite *CategoryServiceSuite) TestDelete() {
	suite.mockDB.EXPECT().DeleteCategory(gomock.Any(), "strategy", false).Return(nil)

	err := suite.service.Delete(context.Background(), "strategy", false)
	suite.Assert().NoError(err)
}

func (suite *CategoryServiceSuite) TestDeleteNotFound() {
	suite.mockDB.EXPECT().DeleteCategory(gomock.Any(), "missing", false).Return(
		middleware.NewError(404, "Record not found"),
	)

	err := suite.service.Delete(context.Background(), "missing", false)
	suite.Assert().Error(err)
}

func TestCategoryServiceSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceSuite))
}
