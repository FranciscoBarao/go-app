package boardgame

import (
	"context"
	"net/http"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/middleware"
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
	return middleware.NewError(http.StatusNotFound, "record not found")
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

func (suite *BoardgameServiceSuite) TestCreate_ExpansionOfExpansionFails() {
	req := testCreateReq("Scenario")
	parentID := uint(2)
	parent := Boardgame{Slug: "seafarers", Name: "Seafarers", MinPlayers: 2, MaxPlayers: 4, BoardgameID: &parentID}

	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "scenario").Return(Boardgame{}, notFound())
	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "seafarers").Return(parent, nil)

	_, err := suite.service.Create(context.Background(), req, "seafarers")
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().ErrorAs(err, &mr)
	suite.Assert().Equal(http.StatusConflict, mr.GetStatus())
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

func (suite *BoardgameServiceSuite) TestGetBySlug() {
	expected := Boardgame{Name: "Catan", Slug: "catan", MinPlayers: 2, MaxPlayers: 4}
	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "catan").Return(expected, nil)

	bg, err := suite.service.GetBySlug(context.Background(), "catan")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, bg)
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

func (suite *BoardgameServiceSuite) TestDeleteBySlug() {
	bg := Boardgame{ID: 1, Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}
	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "catan").Return(bg, nil)
	suite.mockDB.EXPECT().DeleteBoardgame(gomock.Any(), uint(1), false).Return(nil)

	err := suite.service.DeleteBySlug(context.Background(), "catan", false)
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestCreate_InvalidPlayersOrder() {
	req := &CreateBoardgameRequest{Name: "Bad", MinPlayers: 5, MaxPlayers: 2}

	_, err := suite.service.Create(context.Background(), req, "")
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().ErrorAs(err, &mr)
	suite.Assert().Equal(http.StatusBadRequest, mr.GetStatus())
}

func (suite *BoardgameServiceSuite) TestCreate_InvalidPlayTimeOrder() {
	req := &CreateBoardgameRequest{Name: "Bad Time", MinPlayers: 2, MaxPlayers: 4, MinPlayTime: 90, MaxPlayTime: 30}

	_, err := suite.service.Create(context.Background(), req, "")
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().ErrorAs(err, &mr)
	suite.Assert().Equal(http.StatusBadRequest, mr.GetStatus())
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
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().ErrorAs(err, &mr)
	suite.Assert().Equal(http.StatusBadRequest, mr.GetStatus())
}

func (suite *BoardgameServiceSuite) TestCreate_DuplicateSlug() {
	req := testCreateReq("Catan")
	existing := Boardgame{ID: 1, Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}

	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "catan").Return(existing, nil)

	_, err := suite.service.Create(context.Background(), req, "")
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().ErrorAs(err, &mr)
	suite.Assert().Equal(http.StatusConflict, mr.GetStatus())
}

func (suite *BoardgameServiceSuite) TestCreate_NormalizationCollision() {
	req := testCreateReq("Catan!")
	existing := Boardgame{ID: 1, Slug: "catan", Name: "Catan", MinPlayers: 2, MaxPlayers: 4}

	suite.mockDB.EXPECT().GetBoardgameBySlug(gomock.Any(), "catan").Return(existing, nil)

	_, err := suite.service.Create(context.Background(), req, "")
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().ErrorAs(err, &mr)
	suite.Assert().Equal(http.StatusConflict, mr.GetStatus())
}

func (suite *BoardgameServiceSuite) TestCreate_EmptySlug() {
	req := testCreateReq("!!!")

	_, err := suite.service.Create(context.Background(), req, "")
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().ErrorAs(err, &mr)
	suite.Assert().Equal(http.StatusBadRequest, mr.GetStatus())
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

func TestBoardgameServiceSuite(t *testing.T) {
	suite.Run(t, new(BoardgameServiceSuite))
}
