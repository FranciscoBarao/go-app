package mechanism

import (
	"context"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/store"
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
	suite.mockDB.EXPECT().GetMechanismBySlug(gomock.Any(), "worker-placement").Return(Mechanism{}, store.ErrNotFound)
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
	suite.Assert().ErrorIs(err, ErrAlreadyExists)
	suite.Assert().ErrorContains(err, "slug \"worker-placement\" already in use")
}

func (suite *MechanismServiceSuite) TestCreate_EmptySlug() {
	_, err := suite.service.Create(context.Background(), &CreateMechanismRequest{Name: "!!!"})
	suite.Assert().ErrorIs(err, ErrValidation)
	suite.Assert().ErrorContains(err, "name must contain slug-able characters")
}

func (suite *MechanismServiceSuite) TestGet() {
	expected := Mechanism{Slug: "worker-placement", Name: "Worker Placement"}
	suite.mockDB.EXPECT().GetMechanismBySlug(gomock.Any(), "worker-placement").Return(expected, nil)

	m, err := suite.service.Get(context.Background(), "worker-placement")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, m)
}

func (suite *MechanismServiceSuite) TestGet_NotFound() {
	ctx := context.Background()
	suite.mockDB.EXPECT().GetMechanismBySlug(ctx, "missing").Return(Mechanism{}, store.ErrNotFound)

	_, err := suite.service.Get(ctx, "missing")
	suite.Assert().ErrorIs(err, ErrNotFound)
	suite.Assert().ErrorContains(err, "mechanism \"missing\"")
}

func (suite *MechanismServiceSuite) TestGetIDBySlug() {
	ctx := context.Background()
	suite.mockDB.EXPECT().GetMechanismIDBySlug(ctx, "worker-placement").Return(uint(4), nil)

	id, err := suite.service.GetIDBySlug(ctx, "worker-placement")
	suite.Assert().NoError(err)
	suite.Assert().Equal(uint(4), id)
}

func (suite *MechanismServiceSuite) TestGetIDBySlug_NotFound() {
	ctx := context.Background()
	suite.mockDB.EXPECT().GetMechanismIDBySlug(ctx, "missing").Return(uint(0), store.ErrNotFound)

	id, err := suite.service.GetIDBySlug(ctx, "missing")
	suite.Assert().ErrorIs(err, ErrNotFound)
	suite.Assert().Equal(uint(0), id)
}

func (suite *MechanismServiceSuite) TestGetAll() {
	expected := []Mechanism{{Slug: "trading", Name: "Trading"}}
	wantQuery := listopt.Query{
		Pagination: listopt.Pagination{Page: listopt.DefaultPage, PageSize: listopt.DefaultPageSize},
	}
	suite.mockDB.EXPECT().GetAllMechanisms(gomock.Any(), wantQuery).Return(expected, len(expected), nil)

	list, total, err := suite.service.GetAll(context.Background(), listopt.NewQuery())
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, list)
	suite.Assert().Equal(len(expected), total)
}

func (suite *MechanismServiceSuite) TestDelete() {
	suite.mockDB.EXPECT().DeleteMechanism(gomock.Any(), "worker-placement", false).Return(nil)

	err := suite.service.Delete(context.Background(), "worker-placement", false)
	suite.Assert().NoError(err)
}

func (suite *MechanismServiceSuite) TestDelete_Hard() {
	ctx := context.Background()
	suite.mockDB.EXPECT().DeleteMechanism(ctx, "worker-placement", true).Return(nil)

	err := suite.service.Delete(ctx, "worker-placement", true)
	suite.Assert().NoError(err)
}

func (suite *MechanismServiceSuite) TestDeleteNotFound() {
	suite.mockDB.EXPECT().DeleteMechanism(gomock.Any(), "missing", false).Return(store.ErrNotFound)

	err := suite.service.Delete(context.Background(), "missing", false)
	suite.Assert().ErrorIs(err, ErrNotFound)
	suite.Assert().ErrorContains(err, "mechanism \"missing\"")
}

func (suite *MechanismServiceSuite) TestDelete_StillReferenced() {
	ctx := context.Background()
	suite.mockDB.EXPECT().DeleteMechanism(ctx, "worker-placement", false).Return(store.ErrForeignKey)

	err := suite.service.Delete(ctx, "worker-placement", false)
	suite.Assert().ErrorIs(err, ErrStillReferenced)
	suite.Assert().NotErrorIs(err, store.ErrForeignKey)
}

func TestMechanismServiceSuite(t *testing.T) {
	suite.Run(t, new(MechanismServiceSuite))
}
