package boardgame

import (
	"context"
	"net/http"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/middleware"
	tag "github.com/FranciscoBarao/catalog/internal/tag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type BoardgameServiceSuite struct {
	suite.Suite
	mockDB   *MockDatabase
	mockTag  *MockTagService
	mockCat  *MockCategoryService
	mockMech *MockMechanismService
	service  *Service
}

func (suite *BoardgameServiceSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockDB = NewMockDatabase(ctrl)
	suite.mockTag = NewMockTagService(ctrl)
	suite.mockCat = NewMockCategoryService(ctrl)
	suite.mockMech = NewMockMechanismService(ctrl)
	suite.service = NewService(suite.mockDB, suite.mockTag, suite.mockCat, suite.mockMech)
}

func (suite *BoardgameServiceSuite) TestCreate() {
	bg := &Boardgame{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4}
	suite.mockDB.EXPECT().CreateBoardgame(gomock.Any(), bg).Return(nil)

	err := suite.service.Create(context.Background(), bg, 0)
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestCreateWithExpansion() {
	parent := Boardgame{ID: 1, Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4}

	expansion := &Boardgame{Name: "Catan: Seafarers", Publisher: "Kosmos", PlayerNumber: 4}

	suite.mockDB.EXPECT().GetBoardgameByID(gomock.Any(), uint(1)).Return(parent, nil)
	suite.mockDB.EXPECT().CreateBoardgame(gomock.Any(), expansion).Return(nil)

	err := suite.service.Create(context.Background(), expansion, 1)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(expansion.BoardgameID)
	suite.Assert().Equal(uint(1), *expansion.BoardgameID)
}

func (suite *BoardgameServiceSuite) TestCreate_ExpansionOfExpansionFails() {
	parentID := uint(2)
	parent := Boardgame{Name: "Catan: Seafarers", Publisher: "Kosmos", PlayerNumber: 4, BoardgameID: &parentID}

	expansion := &Boardgame{Name: "Catan: Seafarers Scenario", Publisher: "Kosmos", PlayerNumber: 4}

	suite.mockDB.EXPECT().GetBoardgameByID(gomock.Any(), uint(2)).Return(parent, nil)

	err := suite.service.Create(context.Background(), expansion, 2)
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().ErrorAs(err, &mr)
	suite.Assert().Equal(http.StatusConflict, mr.GetStatus())
}

func (suite *BoardgameServiceSuite) TestCreateWithTags() {
	ctx := context.Background()
	bg := &Boardgame{
		Name:         "Catan",
		Publisher:    "Kosmos",
		PlayerNumber: 4,
		Tags:         []tag.Tag{{Name: "strategy"}, {Name: "family"}},
	}

	suite.mockTag.EXPECT().Get(ctx, "strategy").Return(tag.Tag{Name: "strategy"}, nil)
	suite.mockTag.EXPECT().Get(ctx, "family").Return(tag.Tag{Name: "family"}, nil)
	suite.mockDB.EXPECT().CreateBoardgame(ctx, bg).Return(nil)

	err := suite.service.Create(ctx, bg, 0)
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestCreateTag_NotFound() {
	bg := &Boardgame{
		Name:         "Catan",
		Publisher:    "Kosmos",
		PlayerNumber: 4,
		Tags:         []tag.Tag{{Name: "nonexistent"}},
	}

	tagErr := middleware.NewError(http.StatusNotFound, "Tag not found with name: nonexistent")
	suite.mockTag.EXPECT().Get(gomock.Any(), "nonexistent").Return(tag.Tag{}, tagErr)

	err := suite.service.Create(context.Background(), bg, 0)
	suite.Assert().Error(err)
	suite.Assert().Equal("Tag not found with name: nonexistent", err.Error())
}

func (suite *BoardgameServiceSuite) TestCreateWithExpansion_GetIDInternalError() {
	suite.mockDB.EXPECT().
		GetBoardgameByID(gomock.Any(), uint(1)).
		Return(Boardgame{}, assert.AnError)

	err := suite.service.Create(context.Background(), nil, 1)
	suite.Assert().Error(err)
}

func (suite *BoardgameServiceSuite) TestGetByID() {
	parentID := uint(1)
	expected := Boardgame{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4, BoardgameID: &parentID}

	suite.mockDB.EXPECT().GetBoardgameByID(gomock.Any(), uint(1)).Return(expected, nil)

	bg, err := suite.service.GetByID(context.Background(), 1)
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, bg)
}

func (suite *BoardgameServiceSuite) TestGetByID_NotFound() {
	suite.mockDB.EXPECT().GetBoardgameByID(gomock.Any(), uint(999)).Return(
		Boardgame{}, middleware.NewError(http.StatusNotFound, "record not found"),
	)

	_, err := suite.service.GetByID(context.Background(), 999)
	suite.Assert().Error(err)
	suite.Assert().Equal("record not found", err.Error())
}

func (suite *BoardgameServiceSuite) TestUpdate() {
	parentID := uint(1)
	existing := Boardgame{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4, BoardgameID: &parentID}

	name := "Catan Revised"
	playerNumber := 6
	tags := []tag.Tag{{Name: "strategy"}}
	input := &UpdateBoardgameRequest{
		Name:         &name,
		PlayerNumber: &playerNumber,
		Tags:         &tags,
	}

	suite.mockDB.EXPECT().GetBoardgameByID(gomock.Any(), uint(1)).Return(existing, nil)
	suite.mockTag.EXPECT().Get(gomock.Any(), "strategy").Return(tag.Tag{Name: "strategy"}, nil)
	suite.mockDB.EXPECT().UpdateBoardgame(gomock.Any(), gomock.Any()).Return(nil)
	suite.mockDB.EXPECT().ReplaceBoardgameTags(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	err := suite.service.Update(context.Background(), input, 1)
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestDeleteByID() {
	suite.mockDB.EXPECT().DeleteBoardgame(gomock.Any(), uint(1)).Return(nil)

	err := suite.service.DeleteByID(context.Background(), 1)
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestDeleteByID_NotFound() {
	suite.mockDB.EXPECT().DeleteBoardgame(gomock.Any(), uint(999)).Return(
		middleware.NewError(http.StatusNotFound, "record not found"),
	)

	err := suite.service.DeleteByID(context.Background(), 999)
	suite.Assert().Error(err)
	suite.Assert().Equal("record not found", err.Error())
}

func TestBoardgameServiceSuite(t *testing.T) {
	suite.Run(t, new(BoardgameServiceSuite))
}
