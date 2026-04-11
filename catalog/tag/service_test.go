package tag

import (
	"context"
	"testing"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// TagServiceSuite tests the tag Service in isolation using a MockDatabase.
type TagServiceSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	mockDB  *MockDatabase
	service *Service
}

func (suite *TagServiceSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockDB = NewMockDatabase(suite.ctrl)
	suite.service = NewService(suite.mockDB)
}

func (suite *TagServiceSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *TagServiceSuite) TestCreate() {
	tag := NewTag("strategy")
	suite.mockDB.EXPECT().CreateTag(gomock.Any(), tag).Return(nil)

	err := suite.service.Create(context.Background(), tag)
	suite.Assert().NoError(err)
}

func (suite *TagServiceSuite) TestGetAll() {
	expected := []Tag{{Name: "strategy"}, {Name: "cooperative"}}
	suite.mockDB.EXPECT().GetAllTags(gomock.Any(), "name").Return(expected, nil)

	tags, err := suite.service.GetAll(context.Background(), "name")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, tags)
}

func (suite *TagServiceSuite) TestGet() {
	expected := Tag{Name: "strategy"}
	suite.mockDB.EXPECT().GetTag(gomock.Any(), "strategy").Return(expected, nil)

	tag, err := suite.service.Get(context.Background(), "strategy")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, tag)
}

func (suite *TagServiceSuite) TestGetNotFound() {
	suite.mockDB.EXPECT().GetTag(gomock.Any(), "nonexistent").Return(
		Tag{}, middleware.NewError(404, "Record not found"),
	)

	_, err := suite.service.Get(context.Background(), "nonexistent")
	suite.Assert().Error(err)
	suite.Assert().Equal("Record not found", err.Error())
}

func (suite *TagServiceSuite) TestDelete() {
	suite.mockDB.EXPECT().DeleteTag(gomock.Any(), "strategy").Return(nil)

	err := suite.service.Delete(context.Background(), "strategy")
	suite.Assert().NoError(err)
}

func (suite *TagServiceSuite) TestDeleteNotFound() {
	suite.mockDB.EXPECT().DeleteTag(gomock.Any(), "missing").Return(
		middleware.NewError(404, "Record not found"),
	)

	err := suite.service.Delete(context.Background(), "missing")
	suite.Assert().Error(err)
	suite.Assert().Equal("Record not found", err.Error())
}

func TestTagServiceSuite(t *testing.T) {
	suite.Run(t, new(TagServiceSuite))
}
