package contributor

import (
	"context"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/store"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ContributorServiceSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	mockDB  *MockDatabase
	service *Service
}

func (suite *ContributorServiceSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockDB = NewMockDatabase(suite.ctrl)
	suite.service = NewService(suite.mockDB)
}

func (suite *ContributorServiceSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *ContributorServiceSuite) TestCreate() {
	ctx := context.Background()
	bio := "Designer of modern classics"
	bggID := 42
	suite.mockDB.EXPECT().GetContributorBySlug(ctx, "reiner-knizia").Return(Contributor{}, store.ErrNotFound)
	suite.mockDB.EXPECT().CreateContributor(ctx, gomock.AssignableToTypeOf(&Contributor{})).DoAndReturn(
		func(_ context.Context, c *Contributor) error {
			suite.Equal("reiner-knizia", c.Slug)
			suite.Equal("Reiner Knizia", c.Name)
			suite.Require().NotNil(c.Bio)
			suite.Equal(bio, *c.Bio)
			suite.Require().NotNil(c.BggID)
			suite.Equal(bggID, *c.BggID)
			c.ID = 7
			return nil
		},
	)

	got, err := suite.service.Create(ctx, &CreateContributorRequest{Name: "Reiner Knizia", Bio: &bio, BggID: &bggID})
	suite.Assert().NoError(err)
	suite.Assert().Equal(uint(7), got.ID)
	suite.Assert().Equal("reiner-knizia", got.Slug)
}

func (suite *ContributorServiceSuite) TestCreate_DuplicateSlug() {
	ctx := context.Background()
	suite.mockDB.EXPECT().GetContributorBySlug(ctx, "reiner-knizia").Return(Contributor{Slug: "reiner-knizia"}, nil)

	_, err := suite.service.Create(ctx, &CreateContributorRequest{Name: "Reiner Knizia"})
	suite.Assert().ErrorIs(err, ErrAlreadyExists)
	suite.Assert().ErrorContains(err, "slug \"reiner-knizia\" already in use")
}

func (suite *ContributorServiceSuite) TestCreate_EmptySlug() {
	_, err := suite.service.Create(context.Background(), &CreateContributorRequest{Name: "!!!"})
	suite.Assert().ErrorIs(err, ErrValidation)
	suite.Assert().ErrorContains(err, "name must contain slug-able characters")
}

func (suite *ContributorServiceSuite) TestGet() {
	ctx := context.Background()
	expected := Contributor{ID: 7, Slug: "reiner-knizia", Name: "Reiner Knizia"}
	suite.mockDB.EXPECT().GetContributorBySlug(ctx, "reiner-knizia").Return(expected, nil)

	got, err := suite.service.Get(ctx, "reiner-knizia")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, got)
}

func (suite *ContributorServiceSuite) TestGet_NotFound() {
	ctx := context.Background()
	suite.mockDB.EXPECT().GetContributorBySlug(ctx, "missing").Return(Contributor{}, store.ErrNotFound)

	_, err := suite.service.Get(ctx, "missing")
	suite.Assert().ErrorIs(err, ErrNotFound)
	suite.Assert().ErrorContains(err, "contributor \"missing\"")
}

func (suite *ContributorServiceSuite) TestGetIDBySlug() {
	ctx := context.Background()
	suite.mockDB.EXPECT().GetContributorIDBySlug(ctx, "reiner-knizia").Return(uint(7), nil)

	id, err := suite.service.GetIDBySlug(ctx, "reiner-knizia")
	suite.Assert().NoError(err)
	suite.Assert().Equal(uint(7), id)
}

func (suite *ContributorServiceSuite) TestGetIDBySlug_NotFound() {
	ctx := context.Background()
	suite.mockDB.EXPECT().GetContributorIDBySlug(ctx, "missing").Return(uint(0), store.ErrNotFound)

	id, err := suite.service.GetIDBySlug(ctx, "missing")
	suite.Assert().ErrorIs(err, ErrNotFound)
	suite.Assert().Equal(uint(0), id)
}

func (suite *ContributorServiceSuite) TestGetAll() {
	ctx := context.Background()
	expected := []Contributor{{Slug: "reiner-knizia", Name: "Reiner Knizia"}}
	query := listopt.NewQuery()
	suite.mockDB.EXPECT().GetAllContributors(ctx, query).Return(expected, len(expected), nil)

	got, total, err := suite.service.GetAll(ctx, query)
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, got)
	suite.Assert().Equal(len(expected), total)
}

func (suite *ContributorServiceSuite) TestUpdate() {
	ctx := context.Background()
	existing := Contributor{ID: 7, Slug: "reiner-knizia", Name: "Reiner Knizia"}
	name := "Dr. Knizia"
	bio := "updated"
	bggID := 42
	suite.mockDB.EXPECT().GetContributorBySlug(ctx, "reiner-knizia").Return(existing, nil)
	suite.mockDB.EXPECT().UpdateContributor(ctx, gomock.AssignableToTypeOf(&Contributor{})).DoAndReturn(
		func(_ context.Context, c *Contributor) error {
			suite.Equal("Dr. Knizia", c.Name)
			suite.Require().NotNil(c.Bio)
			suite.Equal(bio, *c.Bio)
			suite.Require().NotNil(c.BggID)
			suite.Equal(bggID, *c.BggID)
			suite.Equal("reiner-knizia", c.Slug)
			return nil
		},
	)

	got, err := suite.service.Update(ctx, &UpdateContributorRequest{Name: &name, Bio: &bio, BggID: &bggID}, "reiner-knizia")
	suite.Assert().NoError(err)
	suite.Assert().Equal("Dr. Knizia", got.Name)
	suite.Assert().Equal(&bio, got.Bio)
	suite.Assert().Equal(&bggID, got.BggID)
}

func (suite *ContributorServiceSuite) TestUpdate_EmptyName() {
	name := "   "
	_, err := suite.service.Update(context.Background(), &UpdateContributorRequest{Name: &name}, "reiner-knizia")
	suite.Assert().ErrorIs(err, ErrValidation)
	suite.Assert().ErrorContains(err, "name cannot be empty")
}

func (suite *ContributorServiceSuite) TestUpdate_NotFound() {
	ctx := context.Background()
	name := "Dr. Knizia"
	suite.mockDB.EXPECT().GetContributorBySlug(ctx, "missing").Return(Contributor{}, store.ErrNotFound)

	_, err := suite.service.Update(ctx, &UpdateContributorRequest{Name: &name}, "missing")
	suite.Assert().ErrorIs(err, ErrNotFound)
}

func (suite *ContributorServiceSuite) TestDelete() {
	ctx := context.Background()
	suite.mockDB.EXPECT().DeleteContributor(ctx, "reiner-knizia", true).Return(nil)

	err := suite.service.Delete(ctx, "reiner-knizia", true)
	suite.Assert().NoError(err)
}

func (suite *ContributorServiceSuite) TestDelete_NotFound() {
	ctx := context.Background()
	suite.mockDB.EXPECT().DeleteContributor(ctx, "missing", false).Return(store.ErrNotFound)

	err := suite.service.Delete(ctx, "missing", false)
	suite.Assert().ErrorIs(err, ErrNotFound)
	suite.Assert().ErrorContains(err, "contributor \"missing\"")
}

func (suite *ContributorServiceSuite) TestDelete_StillReferenced() {
	ctx := context.Background()
	suite.mockDB.EXPECT().DeleteContributor(ctx, "reiner-knizia", true).Return(store.ErrForeignKey)

	err := suite.service.Delete(ctx, "reiner-knizia", true)
	suite.Assert().ErrorIs(err, ErrStillReferenced)
	suite.Assert().NotErrorIs(err, store.ErrForeignKey)
}

func TestContributorServiceSuite(t *testing.T) {
	suite.Run(t, new(ContributorServiceSuite))
}
