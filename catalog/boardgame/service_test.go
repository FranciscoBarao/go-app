package boardgame

import (
	"net/http"
	"testing"

	"github.com/FranciscoBarao/catalog/middleware"
	tag "github.com/FranciscoBarao/catalog/tag"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// BoardgameServiceSuite tests the boardgame Service in isolation using mocksuite.
type BoardgameServiceSuite struct {
	suite.Suite
	mockDB   *MockDatabase
	mockTag  *MockTagGetter
	mockCat  *MockCategoryGetter
	mockMech *MockMechanismGetter
	service  *Service
}

func (suite *BoardgameServiceSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockDB = NewMockDatabase(ctrl)
	suite.mockTag = NewMockTagGetter(ctrl)
	suite.mockCat = NewMockCategoryGetter(ctrl)
	suite.mockMech = NewMockMechanismGetter(ctrl)
	suite.service = NewService(suite.mockDB, suite.mockTag, suite.mockCat, suite.mockMech)
}

func (suite *BoardgameServiceSuite) TestCreate() {
	bg := &Boardgame{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4}
	suite.mockDB.EXPECT().Create(bg).Return(nil)

	err := suite.service.Create(bg, "")
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestCreateWithExpansion() {
	parentID := uint(1)
	parent := Boardgame{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4, BoardgameID: &parentID}

	expansion := &Boardgame{Name: "Catan: Seafarers", Publisher: "Kosmos", PlayerNumber: 4}

	// GetByID for parent — db.Read populates the dest pointer
	suite.mockDB.EXPECT().Read(gomock.Any(), "", "id = ?", "1").DoAndReturn(
		func(dest interface{}, sort, search, identifier string) error {
			ptr := dest.(*Boardgame)
			*ptr = parent
			return nil
		},
	)
	suite.mockDB.EXPECT().Create(expansion).Return(nil)

	err := suite.service.Create(expansion, "1")
	suite.Assert().NoError(err)
	suite.Assert().NotNil(expansion.BoardgameID)
	suite.Assert().Equal(parentID, *expansion.BoardgameID)
}

func (suite *BoardgameServiceSuite) TestCreateExpansionOfExpansionFails() {
	parentID := uint(2)
	parent := Boardgame{Name: "Catan: Seafarers", Publisher: "Kosmos", PlayerNumber: 4, BoardgameID: &parentID}

	expansion := &Boardgame{Name: "Catan: Seafarers Scenario", Publisher: "Kosmos", PlayerNumber: 4}

	// GetByID for parent — parent is an expansion (BoardgameID != nil)
	suite.mockDB.EXPECT().Read(gomock.Any(), "", "id = ?", "2").DoAndReturn(
		func(dest interface{}, sort, search, identifier string) error {
			ptr := dest.(*Boardgame)
			*ptr = parent
			return nil
		},
	)
	// db.Create should NOT be called

	err := suite.service.Create(expansion, "2")
	suite.Assert().Error(err)

	var mr *middleware.MalformedRequest
	suite.Assert().ErrorAs(err, &mr)
	suite.Assert().Equal(http.StatusConflict, mr.GetStatus())
}

func (suite *BoardgameServiceSuite) TestCreateWithTags() {
	bg := &Boardgame{
		Name:         "Catan",
		Publisher:    "Kosmos",
		PlayerNumber: 4,
		Tags:         []tag.Tag{{Name: "strategy"}, {Name: "family"}},
	}

	suite.mockTag.EXPECT().Get("strategy").Return(tag.Tag{Name: "strategy"}, nil)
	suite.mockTag.EXPECT().Get("family").Return(tag.Tag{Name: "family"}, nil)
	suite.mockDB.EXPECT().Create(bg).Return(nil)

	err := suite.service.Create(bg, "")
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestCreateTagNotFound() {
	bg := &Boardgame{
		Name:         "Catan",
		Publisher:    "Kosmos",
		PlayerNumber: 4,
		Tags:         []tag.Tag{{Name: "nonexistent"}},
	}

	tagErr := middleware.NewError(http.StatusNotFound, "Tag not found with name: nonexistent")
	suite.mockTag.EXPECT().Get("nonexistent").Return(tag.Tag{}, tagErr)
	// db.Create should NOT be called

	err := suite.service.Create(bg, "")
	suite.Assert().Error(err)
	suite.Assert().Equal("Tag not found with name: nonexistent", err.Error())
}

func (suite *BoardgameServiceSuite) TestGetByID() {
	parentID := uint(1)
	expected := Boardgame{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4, BoardgameID: &parentID}

	suite.mockDB.EXPECT().Read(gomock.Any(), "", "id = ?", "1").DoAndReturn(
		func(dest interface{}, sort, search, identifier string) error {
			ptr := dest.(*Boardgame)
			*ptr = expected
			return nil
		},
	)

	bg, err := suite.service.GetByID("1")
	suite.Assert().NoError(err)
	suite.Assert().Equal(expected, bg)
}

func (suite *BoardgameServiceSuite) TestGetByIDNotFound() {
	suite.mockDB.EXPECT().Read(gomock.Any(), "", "id = ?", "999").Return(
		middleware.NewError(http.StatusNotFound, "not found"),
	)

	_, err := suite.service.GetByID("999")
	suite.Assert().Error(err)
	suite.Assert().Equal("Boardgame not found with id: 999", err.Error())
}

func (suite *BoardgameServiceSuite) TestUpdate() {
	parentID := uint(1)
	existing := Boardgame{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4, BoardgameID: &parentID}

	input := &Boardgame{
		Name:         "Catan Revised",
		Publisher:    "Kosmos",
		PlayerNumber: 6,
		Tags:         []tag.Tag{{Name: "strategy"}},
	}

	// validateAssociations — tag lookup
	suite.mockTag.EXPECT().Get("strategy").Return(tag.Tag{Name: "strategy"}, nil)

	// GetByID — fetch existing
	suite.mockDB.EXPECT().Read(gomock.Any(), "", "id = ?", "1").DoAndReturn(
		func(dest interface{}, sort, search, identifier string) error {
			ptr := dest.(*Boardgame)
			*ptr = existing
			return nil
		},
	)

	// After UpdateBoardgame is applied, the existing boardgame gets input's fields
	suite.mockDB.EXPECT().Update(gomock.Any()).Return(nil)
	suite.mockDB.EXPECT().ReplaceAssociatons(gomock.Any(), "Tags", gomock.Any()).Return(nil)

	err := suite.service.Update(input, "1")
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestDeleteByID() {
	parentID := uint(1)
	existing := Boardgame{Name: "Catan", Publisher: "Kosmos", PlayerNumber: 4, BoardgameID: &parentID}

	suite.mockDB.EXPECT().Read(gomock.Any(), "", "id = ?", "1").DoAndReturn(
		func(dest interface{}, sort, search, identifier string) error {
			ptr := dest.(*Boardgame)
			*ptr = existing
			return nil
		},
	)
	suite.mockDB.EXPECT().Delete(gomock.Any()).Return(nil)

	err := suite.service.DeleteByID("1")
	suite.Assert().NoError(err)
}

func (suite *BoardgameServiceSuite) TestDeleteByIDNotFound() {
	suite.mockDB.EXPECT().Read(gomock.Any(), "", "id = ?", "999").Return(
		middleware.NewError(http.StatusNotFound, "not found"),
	)
	// db.Delete should NOT be called

	err := suite.service.DeleteByID("999")
	suite.Assert().Error(err)
	suite.Assert().Equal("Boardgame not found with id: 999", err.Error())
}

func TestBoardgameServiceSuite(t *testing.T) {
	suite.Run(t, new(BoardgameServiceSuite))
}
