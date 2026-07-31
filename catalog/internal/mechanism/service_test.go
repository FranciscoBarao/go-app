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
	suite.mockDB.EXPECT().GetMechanismBySlug(gomock.Any(), "worker-placement").Return(Mechanism{}, middleware.NewError(404, "record not found"))
	suite.mockDB.EXPECT().CreateMechanism(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, m *Mechanism) error {
			suite.Equal("worker-placement", m.Slug)
			suite.Equal("Worker Placement", m.Name)
			return nil
		},
	)

	m, err := suite.service.Create(context.Background(), &CreateMechanismRequest{Name: "Worker Placement"})
	suite.Assert().NoError(err)
	suite.Assert().Equal("worker-placement", m.Slug)
}

func (suite *MechanismServiceSuite) TestCreate_DuplicateSlug() {
	suite.mockDB.EXPECT().GetMechanismBySlug(gomock.Any(), "worker-placement").Return(Mechanism{Slug: "worker-placement", Name: "Worker Placement"}, nil)

	_, err := suite.service.Create(context.Background(), &CreateMechanismRequest{Name: "Worker Placement"})
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().ErrorAs(err, &mr)
	suite.Assert().Equal(409, mr.GetStatus())
}

func (suite *MechanismServiceSuite) TestCreate_EmptySlug() {
	_, err := suite.service.Create(context.Background(), &CreateMechanismRequest{Name: "!!!"})
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().ErrorAs(err, &mr)
	suite.Assert().Equal(400, mr.GetStatus())
}

func (suite *MechanismServiceSuite) TestGet() {
	expected := Mechanism{Slug: "worker-placement", Name: "Worker Placement"}
	suite.mockDB.EXPECT().GetMechanismBySlug(gomock.Any(), "worker-placement").Return(expected, nil)

	m, err := suite.service.Get(context.Background(), "worker-placement")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, m)
}

func (suite *MechanismServiceSuite) TestDelete() {
	suite.mockDB.EXPECT().DeleteMechanism(gomock.Any(), "worker-placement", false).Return(nil)

	err := suite.service.Delete(context.Background(), "worker-placement", false)
	suite.Assert().NoError(err)
}

func (suite *MechanismServiceSuite) TestDeleteNotFound() {
	suite.mockDB.EXPECT().DeleteMechanism(gomock.Any(), "missing", false).Return(
		middleware.NewError(404, "Record not found"),
	)

	err := suite.service.Delete(context.Background(), "missing", false)
	suite.Assert().Error(err)
}

func (suite *MechanismServiceSuite) TestGetAll() {
	expected := []Mechanism{{Slug: "trading", Name: "Trading"}}
	wantParams := listopt.Params{
		Pagination: listopt.Pagination{Page: listopt.DefaultPage, PageSize: listopt.DefaultPageSize},
	}
	suite.mockDB.EXPECT().GetAllMechanisms(gomock.Any(), wantParams).Return(expected, len(expected), nil)

	list, total, err := suite.service.GetAll(context.Background())
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, list)
	suite.Assert().Equal(len(expected), total)
}

func TestMechanismServiceSuite(t *testing.T) {
	suite.Run(t, new(MechanismServiceSuite))
}
