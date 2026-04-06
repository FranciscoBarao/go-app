package tag

import (
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
	suite.mockDB.EXPECT().Create(tag).Return(nil)

	err := suite.service.Create(tag)
	suite.Assert().NoError(err)
}

func (suite *TagServiceSuite) TestGetAll() {
	expected := []Tag{{Name: "strategy"}, {Name: "cooperative"}}
	suite.mockDB.EXPECT().Read(gomock.Any(), "name", "", "").DoAndReturn(
		func(dest interface{}, sort, search, identifier string) error {
			ptr := dest.(*[]Tag)
			*ptr = expected
			return nil
		},
	)

	tags, err := suite.service.GetAll("name")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, tags)
}

func (suite *TagServiceSuite) TestGet() {
	expected := Tag{Name: "strategy"}
	suite.mockDB.EXPECT().Read(gomock.Any(), "", "name = ?", "strategy").DoAndReturn(
		func(dest interface{}, sort, search, identifier string) error {
			ptr := dest.(*Tag)
			*ptr = expected
			return nil
		},
	)

	tag, err := suite.service.Get("strategy")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, tag)
}

func (suite *TagServiceSuite) TestGetNotFound() {
	suite.mockDB.EXPECT().Read(gomock.Any(), "", "name = ?", "nonexistent").Return(
		middleware.NewError(404, "not found"),
	)

	_, err := suite.service.Get("nonexistent")
	suite.Assert().Error(err)
	suite.Assert().Equal("Tag not found with name: nonexistent", err.Error())
}

func (suite *TagServiceSuite) TestDelete() {
	expected := Tag{Name: "strategy"}
	suite.mockDB.EXPECT().Read(gomock.Any(), "", "name = ?", "strategy").DoAndReturn(
		func(dest interface{}, sort, search, identifier string) error {
			ptr := dest.(*Tag)
			*ptr = expected
			return nil
		},
	)
	suite.mockDB.EXPECT().Delete(&expected).Return(nil)

	err := suite.service.Delete("strategy")
	suite.Assert().NoError(err)
}

func (suite *TagServiceSuite) TestDeleteNotFound() {
	suite.mockDB.EXPECT().Read(gomock.Any(), "", "name = ?", "missing").Return(
		middleware.NewError(404, "not found"),
	)
	// db.Delete should NOT be called

	err := suite.service.Delete("missing")
	suite.Assert().Error(err)
	suite.Assert().Equal("Tag not found with name: missing", err.Error())
}

func TestTagServiceSuite(t *testing.T) {
	suite.Run(t, new(TagServiceSuite))
}
