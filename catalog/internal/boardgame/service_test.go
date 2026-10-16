package boardgame

import (
	"context"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type BoardgameServiceSuite struct {
	suite.Suite
	mockDB      *MockDatabase
	mockCat     *MockCategoryService
	mockMech    *MockMechanismService
	mockContrib *MockContributorService
	service     *Service
}

func (suite *BoardgameServiceSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockDB = NewMockDatabase(ctrl)
	suite.mockCat = NewMockCategoryService(ctrl)
	suite.mockMech = NewMockMechanismService(ctrl)
	suite.mockContrib = NewMockContributorService(ctrl)
	suite.service = NewService(suite.mockDB, suite.mockCat, suite.mockMech, suite.mockContrib)
}

func testCreateReq(name string) *CreateBoardgameRequest {
	return &CreateBoardgameRequest{Name: name, MinPlayers: 2, MaxPlayers: 4}
}

func notFound() error {
	return store.ErrNotFound
}

func (suite *BoardgameServiceSuite) TestCreate() {
	req := testCreateReq("Catan")
	expected := Boardgame{ID: 1, Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}

	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "catan").Return(Boardgame{}, notFound())
	suite.mockDB.EXPECT().CreateBoardgame(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, input CreateBoardgameDTO) (Boardgame, error) {
			suite.Equal("catan", input.Slug)
			suite.Equal("Catan", input.Name)
			return expected, nil
		},
	)

	bg, err := suite.service.Create(context.Background(), req, "")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, bg)
}

func (suite *BoardgameServiceSuite) TestCreateWithExpansion() {
	req := testCreateReq("Catan: Seafarers")
	parent := Boardgame{ID: 1, Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}
	parentID := uint(1)
	expected := Boardgame{ID: 2, Slug: "catan-seafarers", Name: "Catan: Seafarers", MinPlayers: 2, MaxPlayers: 4, BoardgameID: &parentID}

	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "catan-seafarers").Return(Boardgame{}, notFound())
	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "catan").Return(parent, nil)
	suite.mockDB.EXPECT().CreateBoardgame(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, input CreateBoardgameDTO) (Boardgame, error) {
			suite.Require().NotNil(input.ParentID)
			suite.Equal(uint(1), *input.ParentID)
			return expected, nil
		},
	)

	bg, err := suite.service.Create(context.Background(), req, "catan")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, bg)
}

func (suite *BoardgameServiceSuite) TestCreateWithCategories() {
	ctx := context.Background()
	req := &CreateBoardgameRequest{
		Name:       "Catan",
		MinPlayers: 2,
		MaxPlayers: 4,
		Categories: []CategoryRef{{Slug: "economic"}},
	}

	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "catan").Return(Boardgame{}, notFound())
	suite.mockCat.EXPECT().GetIDBySlug(ctx, "economic").Return(uint(1), nil)
	suite.mockDB.EXPECT().CreateBoardgame(ctx, gomock.Any()).DoAndReturn(
		func(_ context.Context, input CreateBoardgameDTO) (Boardgame, error) {
			suite.Len(input.Categories, 1)
			suite.Equal(uint(1), input.Categories[0])
			return Boardgame{Slug: "catan", Name: "Catan"}, nil
		},
	)

	_, err := suite.service.Create(ctx, req, "")
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestCreate_PropagatesScalarFields() {
	req := &CreateBoardgameRequest{
		Name:          "Catan",
		MinPlayers:    2,
		MaxPlayers:    4,
		Description:   "A trading game",
		YearPublished: 1995,
		MinPlayTime:   45,
		MaxPlayTime:   90,
	}

	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "catan").Return(Boardgame{}, notFound())
	suite.mockDB.EXPECT().CreateBoardgame(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, input CreateBoardgameDTO) (Boardgame, error) {
			suite.Equal("A trading game", input.Description)
			suite.Equal(1995, input.YearPublished)
			suite.Equal(45, input.MinPlayTime)
			suite.Equal(90, input.MaxPlayTime)
			return Boardgame{Slug: "catan", Name: "Catan"}, nil
		},
	)

	_, err := suite.service.Create(context.Background(), req, "")
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestCreate_ExpansionOfExpansionFails() {
	req := testCreateReq("Scenario")
	parentID := uint(2)
	parent := Boardgame{Slug: "seafarers", Name: "Seafarers", MinPlayers: 2, MaxPlayers: 4, BoardgameID: &parentID}

	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "scenario").Return(Boardgame{}, notFound())
	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "seafarers").Return(parent, nil)

	_, err := suite.service.Create(context.Background(), req, "seafarers")
	suite.Assert().ErrorIs(err, ErrParentIsExpansion)
	suite.Assert().ErrorContains(err, "expansion \"seafarers\": an expansion cannot be a parent")
}

func (suite *BoardgameServiceSuite) TestCreate_InvalidPlayersOrder() {
	req := &CreateBoardgameRequest{Name: "Bad", MinPlayers: 5, MaxPlayers: 2}

	_, err := suite.service.Create(context.Background(), req, "")
	suite.Assert().ErrorIs(err, ErrValidation)
	suite.Assert().ErrorContains(err, "min_players cannot exceed max_players")
}

func (suite *BoardgameServiceSuite) TestCreate_InvalidPlayTimeOrder() {
	req := &CreateBoardgameRequest{Name: "Bad Time", MinPlayers: 2, MaxPlayers: 4, MinPlayTime: 90, MaxPlayTime: 30}

	_, err := suite.service.Create(context.Background(), req, "")
	suite.Assert().ErrorIs(err, ErrValidation)
	suite.Assert().ErrorContains(err, "min_play_time cannot exceed max_play_time")
}

func (suite *BoardgameServiceSuite) TestCreate_InvalidContributionRole() {
	ctx := context.Background()
	req := &CreateBoardgameRequest{
		Name:       "Catan",
		MinPlayers: 2,
		MaxPlayers: 4,
		Contributions: []contributor.ContributionInput{{
			Slug: "kosmos",
			Role: contributor.Role("invalid"),
		}},
	}

	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "catan").Return(Boardgame{}, notFound())

	_, err := suite.service.Create(ctx, req, "")
	suite.Assert().ErrorIs(err, ErrValidation)
	suite.Assert().ErrorContains(err, "contribution role \"invalid\"")
}

func (suite *BoardgameServiceSuite) TestCreate_DuplicateSlug() {
	req := testCreateReq("Catan")
	existing := Boardgame{ID: 1, Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}

	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "catan").Return(existing, nil)

	_, err := suite.service.Create(context.Background(), req, "")
	suite.Assert().ErrorIs(err, ErrAlreadyExists)
	suite.Assert().ErrorContains(err, "slug \"catan\" already in use")
}

func (suite *BoardgameServiceSuite) TestCreate_NormalizationCollision() {
	req := testCreateReq("Catan!")
	existing := Boardgame{ID: 1, Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}

	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "catan").Return(existing, nil)

	_, err := suite.service.Create(context.Background(), req, "")
	suite.Assert().ErrorIs(err, ErrAlreadyExists)
	suite.Assert().ErrorContains(err, "slug \"catan\" already in use")
}

func (suite *BoardgameServiceSuite) TestCreate_EmptySlug() {
	req := testCreateReq("!!!")

	_, err := suite.service.Create(context.Background(), req, "")
	suite.Assert().ErrorIs(err, ErrValidation)
	suite.Assert().ErrorContains(err, "name must contain slug-able characters")
}

func (suite *BoardgameServiceSuite) TestCreate_MissingCategory() {
	ctx := context.Background()
	req := testCreateReq("Catan")
	req.Categories = []CategoryRef{{Slug: "missing"}}
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "catan").Return(Boardgame{}, notFound())
	suite.mockCat.EXPECT().GetIDBySlug(ctx, "missing").Return(uint(0), category.ErrNotFound)

	_, err := suite.service.Create(ctx, req, "")
	suite.Assert().ErrorIs(err, ErrNotFound)
	suite.Assert().ErrorContains(err, "category \"missing\"")
}

func (suite *BoardgameServiceSuite) TestCreate_CategoryLookupError() {
	ctx := context.Background()
	req := testCreateReq("Catan")
	req.Categories = []CategoryRef{{Slug: "strategy"}}
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "catan").Return(Boardgame{}, notFound())
	suite.mockCat.EXPECT().GetIDBySlug(ctx, "strategy").Return(uint(0), assert.AnError)

	_, err := suite.service.Create(ctx, req, "")
	suite.Assert().ErrorContains(err, "resolve category \"strategy\": internal error")
	suite.Assert().NotErrorIs(err, assert.AnError)
}

func (suite *BoardgameServiceSuite) TestCreate_ParentNotFound() {
	ctx := context.Background()
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "scenario").Return(Boardgame{}, notFound())
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "missing").Return(Boardgame{}, store.ErrNotFound)

	_, err := suite.service.Create(ctx, testCreateReq("Scenario"), "missing")
	suite.Assert().ErrorIs(err, ErrNotFound)
	suite.Assert().ErrorContains(err, "boardgame \"missing\"")
}

func (suite *BoardgameServiceSuite) TestCreate_MissingMechanism() {
	ctx := context.Background()
	req := testCreateReq("Catan")
	req.Mechanisms = []MechanismRef{{Slug: "missing"}}
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "catan").Return(Boardgame{}, notFound())
	suite.mockMech.EXPECT().GetIDBySlug(ctx, "missing").Return(uint(0), mechanism.ErrNotFound)

	_, err := suite.service.Create(ctx, req, "")
	suite.Assert().ErrorIs(err, ErrNotFound)
	suite.Assert().ErrorContains(err, "mechanism \"missing\"")
}

func (suite *BoardgameServiceSuite) TestCreate_MechanismLookupError() {
	ctx := context.Background()
	req := testCreateReq("Catan")
	req.Mechanisms = []MechanismRef{{Slug: "auction"}}
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "catan").Return(Boardgame{}, notFound())
	suite.mockMech.EXPECT().GetIDBySlug(ctx, "auction").Return(uint(0), assert.AnError)

	_, err := suite.service.Create(ctx, req, "")
	suite.Assert().ErrorContains(err, "resolve mechanism \"auction\": internal error")
	suite.Assert().NotErrorIs(err, assert.AnError)
}

func (suite *BoardgameServiceSuite) TestCreate_MissingContributor() {
	ctx := context.Background()
	req := testCreateReq("Catan")
	req.Contributions = []contributor.ContributionInput{{Slug: "missing", Role: contributor.RoleDesigner}}
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "catan").Return(Boardgame{}, notFound())
	suite.mockContrib.EXPECT().GetIDBySlug(ctx, "missing").Return(uint(0), contributor.ErrNotFound)

	_, err := suite.service.Create(ctx, req, "")
	suite.Assert().ErrorIs(err, ErrNotFound)
	suite.Assert().ErrorContains(err, "contributor \"missing\"")
}

func (suite *BoardgameServiceSuite) TestCreate_ContributorLookupError() {
	ctx := context.Background()
	req := testCreateReq("Catan")
	req.Contributions = []contributor.ContributionInput{{Slug: "kosmos", Role: contributor.RoleDesigner}}
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "catan").Return(Boardgame{}, notFound())
	suite.mockContrib.EXPECT().GetIDBySlug(ctx, "kosmos").Return(uint(0), assert.AnError)

	_, err := suite.service.Create(ctx, req, "")
	suite.Assert().ErrorContains(err, "resolve contributor \"kosmos\": internal error")
	suite.Assert().NotErrorIs(err, assert.AnError)
}

func (suite *BoardgameServiceSuite) TestGetBySlug() {
	expected := Boardgame{Name: "Catan", Slug: "catan", MinPlayers: 2, MaxPlayers: 4}
	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "catan").Return(expected, nil)

	bg, err := suite.service.GetBySlug(context.Background(), "catan")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, bg)
}

func (suite *BoardgameServiceSuite) TestGetBySlug_NotFound() {
	ctx := context.Background()
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "missing").Return(Boardgame{}, store.ErrNotFound)

	_, err := suite.service.GetBySlug(ctx, "missing")
	suite.Assert().ErrorIs(err, ErrNotFound)
	suite.Assert().ErrorContains(err, "boardgame \"missing\"")
}

func (suite *BoardgameServiceSuite) TestGetByID() {
	ctx := context.Background()
	expected := Boardgame{ID: 9, Slug: "catan", Name: "Catan"}
	suite.mockDB.EXPECT().GetBoardgameByID(ctx, uint(9)).Return(expected, nil)

	got, err := suite.service.GetByID(ctx, 9)
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, got)
}

func (suite *BoardgameServiceSuite) TestGetByID_NotFound() {
	ctx := context.Background()
	suite.mockDB.EXPECT().GetBoardgameByID(ctx, uint(9)).Return(Boardgame{}, store.ErrNotFound)

	_, err := suite.service.GetByID(ctx, 9)
	suite.Assert().ErrorIs(err, ErrNotFound)
	suite.Assert().ErrorContains(err, "boardgame 9")
}

func (suite *BoardgameServiceSuite) TestGetAll() {
	expected := []Boardgame{{Slug: "catan", Name: "Catan"}}
	wantQuery := listopt.Query{
		Sort:       listopt.Sort{Column: "name", Order: "asc"},
		Pagination: listopt.Pagination{Page: listopt.DefaultPage, PageSize: listopt.DefaultPageSize},
	}
	suite.mockDB.EXPECT().
		GetAllBoardgames(gomock.Any(), wantQuery, true).
		Return(expected, len(expected), nil)

	all, total, err := suite.service.GetAll(context.Background(), true, listopt.NewQuery(listopt.WithSort(listopt.Sort{Column: "name", Order: "asc"})))
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, all)
	suite.Assert().Equal(len(expected), total)
}

func (suite *BoardgameServiceSuite) TestUpdate() {
	existing := Boardgame{ID: 1, Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}

	name := "Catan Revised"
	maxPlayers := 6
	mechRefs := []MechanismRef{{Slug: "trading"}}
	input := &UpdateBoardgameRequest{
		Name:       &name,
		MaxPlayers: &maxPlayers,
		Mechanisms: &mechRefs,
	}

	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "catan").Return(existing, nil)
	suite.mockMech.EXPECT().GetIDBySlug(gomock.Any(), "trading").Return(uint(1), nil)
	suite.mockDB.EXPECT().UpdateBoardgameWithAssociations(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	err := suite.service.Update(context.Background(), input, "catan")
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestUpdate_ReplacesCategories() {
	ctx := context.Background()
	existing := Boardgame{ID: 1, Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}
	description := "A trading game"
	refs := []CategoryRef{{Slug: "strategy"}}
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "catan").Return(existing, nil)
	suite.mockCat.EXPECT().GetIDBySlug(ctx, "strategy").Return(uint(3), nil)
	suite.mockDB.EXPECT().UpdateBoardgameWithAssociations(ctx, gomock.AssignableToTypeOf(&Boardgame{}), gomock.AssignableToTypeOf(UpdateAssociations{})).DoAndReturn(
		func(_ context.Context, bg *Boardgame, assoc UpdateAssociations) error {
			suite.Equal(description, bg.Description)
			suite.Equal("catan", bg.Slug)
			suite.Require().NotNil(assoc.Categories)
			suite.Equal([]uint{3}, *assoc.Categories)
			suite.Nil(assoc.Mechanisms)
			suite.Nil(assoc.Contributions)
			return nil
		},
	)

	err := suite.service.Update(ctx, &UpdateBoardgameRequest{Description: &description, Categories: &refs}, "catan")
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestUpdate_NotFound() {
	ctx := context.Background()
	description := "A trading game"
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "missing").Return(Boardgame{}, store.ErrNotFound)

	err := suite.service.Update(ctx, &UpdateBoardgameRequest{Description: &description}, "missing")
	suite.Assert().ErrorIs(err, ErrNotFound)
}

func (suite *BoardgameServiceSuite) TestUpdate_InvalidPlayersOrder() {
	ctx := context.Background()
	existing := Boardgame{ID: 1, Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}
	maxPlayers := 1
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "catan").Return(existing, nil)

	err := suite.service.Update(ctx, &UpdateBoardgameRequest{MaxPlayers: &maxPlayers}, "catan")
	suite.Assert().ErrorIs(err, ErrValidation)
	suite.Assert().ErrorContains(err, "min_players cannot exceed max_players")
}

func (suite *BoardgameServiceSuite) TestUpdate_InvalidPlayTimeOrder() {
	ctx := context.Background()
	existing := Boardgame{ID: 1, Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4, MinPlayTime: 30, MaxPlayTime: 60}
	maxPlayTime := 10
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "catan").Return(existing, nil)

	err := suite.service.Update(ctx, &UpdateBoardgameRequest{MaxPlayTime: &maxPlayTime}, "catan")
	suite.Assert().ErrorIs(err, ErrValidation)
	suite.Assert().ErrorContains(err, "min_play_time cannot exceed max_play_time")
}

func (suite *BoardgameServiceSuite) TestUpdate_MissingCategory() {
	ctx := context.Background()
	existing := Boardgame{ID: 1, Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}
	refs := []CategoryRef{{Slug: "missing"}}
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "catan").Return(existing, nil)
	suite.mockCat.EXPECT().GetIDBySlug(ctx, "missing").Return(uint(0), category.ErrNotFound)

	err := suite.service.Update(ctx, &UpdateBoardgameRequest{Categories: &refs}, "catan")
	suite.Assert().ErrorIs(err, ErrNotFound)
	suite.Assert().ErrorContains(err, "category \"missing\"")
}

func (suite *BoardgameServiceSuite) TestDeleteBySlug() {
	bg := Boardgame{ID: 1, Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}
	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "catan").Return(bg, nil)
	suite.mockDB.EXPECT().DeleteBoardgame(gomock.Any(), uint(1), false).Return(nil)

	err := suite.service.DeleteBySlug(context.Background(), "catan", false)
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestDeleteBySlug_NotFound() {
	ctx := context.Background()
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "missing").Return(Boardgame{}, store.ErrNotFound)

	err := suite.service.DeleteBySlug(ctx, "missing", false)
	suite.Assert().ErrorIs(err, ErrNotFound)
}

func (suite *BoardgameServiceSuite) TestRate() {
	ctx := context.Background()
	rating := &Rating{Value: 8}
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "catan").Return(Boardgame{ID: 1, Slug: "catan"}, nil)

	err := suite.service.Rate(ctx, rating, "catan", "ada")
	suite.Assert().NoError(err)
	suite.Assert().Equal("ada", rating.Username)
}

func (suite *BoardgameServiceSuite) TestRate_NotFound() {
	ctx := context.Background()
	suite.mockDB.EXPECT().GetBoardgameBySlug(ctx, "missing").Return(Boardgame{}, store.ErrNotFound)

	err := suite.service.Rate(ctx, &Rating{Value: 8}, "missing", "ada")
	suite.Assert().ErrorIs(err, ErrNotFound)
}

func TestBoardgameServiceSuite(t *testing.T) {
	suite.Run(t, new(BoardgameServiceSuite))
}
