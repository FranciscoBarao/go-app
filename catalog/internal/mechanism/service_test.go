package mechanism

import (
	"context"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type MechanismServiceSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	mockDB  *MockDatabase
	service *Service
}

func (suite *MechanismServiceSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockDB = NewMockDatabase(suite.ctrl)
	suite.service = NewService(suite.mockDB)
}

func (suite *MechanismServiceSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *MechanismServiceSuite) TestCreate() {
	m := NewMechanism("deckbuilding")
	suite.mockDB.EXPECT().CreateMechanism(gomock.Any(), m).Return(nil)

	err := suite.service.Create(context.Background(), m)
	suite.Assert().NoError(err)
}

func (suite *MechanismServiceSuite) TestGetAll() {
	expected := []Mechanism{{Name: "deckbuilding"}, {Name: "workerplacement"}}
	suite.mockDB.EXPECT().GetAllMechanisms(gomock.Any(), listopt.Params{Sort: listopt.Sort{Column: "name", Order: "asc"}}).Return(expected, nil)

	mechanisms, err := suite.service.GetAll(context.Background(), listopt.WithSort("name", "asc"))
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, mechanisms)
}

func (suite *MechanismServiceSuite) TestGetAllNoSort() {
	expected := []Mechanism{{Name: "deckbuilding"}}
	suite.mockDB.EXPECT().GetAllMechanisms(gomock.Any(), listopt.Params{}).Return(expected, nil)

	mechanisms, err := suite.service.GetAll(context.Background())
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, mechanisms)
}

func (suite *MechanismServiceSuite) TestGet() {
	expected := Mechanism{Name: "deckbuilding"}
	suite.mockDB.EXPECT().GetMechanism(gomock.Any(), "deckbuilding").Return(expected, nil)

	m, err := suite.service.Get(context.Background(), "deckbuilding")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, m)
}

func (suite *MechanismServiceSuite) TestGetNotFound() {
	suite.mockDB.EXPECT().GetMechanism(gomock.Any(), "nonexistent").Return(
		Mechanism{}, middleware.NewError(404, "Record not found"),
	)

	_, err := suite.service.Get(context.Background(), "nonexistent")
	suite.Assert().Error(err)
	suite.Assert().Equal("Record not found", err.Error())
}

func (suite *MechanismServiceSuite) TestDelete() {
	suite.mockDB.EXPECT().DeleteMechanism(gomock.Any(), "deckbuilding").Return(nil)

	err := suite.service.Delete(context.Background(), "deckbuilding")
	suite.Assert().NoError(err)
}

func (suite *MechanismServiceSuite) TestDeleteNotFound() {
	suite.mockDB.EXPECT().DeleteMechanism(gomock.Any(), "missing").Return(
		middleware.NewError(404, "Record not found"),
	)

	err := suite.service.Delete(context.Background(), "missing")
	suite.Assert().Error(err)
	suite.Assert().Equal("Record not found", err.Error())
}

func TestMechanismServiceSuite(t *testing.T) {
	suite.Run(t, new(MechanismServiceSuite))
}
